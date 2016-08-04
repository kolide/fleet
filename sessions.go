package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

var (
	// An error returned by SessionBackend.Get() if no session record was found
	// in the database
	ErrNoActiveSession = errors.New("Active session is not present in the database")

	// An error returned by SessionBackend.Session() if no session has been
	// created yet
	ErrSessionNotCreated = errors.New("The session has not been created")

	ErrSessionExpired = errors.New("The session has expired")
)

const (
	// The name of the session cookie
	CookieName = "KolideSession"
)

// Session is the model object which represents what an active session is
type Session struct {
	BaseModel
	UserID     uint   `gorm:"not null"`
	Key        string `gorm:"not null;unique_index:idx_session_unique_key"`
	AccessedAt time.Time
}

////////////////////////////////////////////////////////////////////////////////
// Managing sessions
////////////////////////////////////////////////////////////////////////////////

type SessionManager struct {
	backend SessionBackend
	request *http.Request
	writer  http.ResponseWriter
	session *Session
	db      *gorm.DB
}

func NewSessionManager(request *http.Request, writer http.ResponseWriter, backend SessionBackend, db *gorm.DB) *SessionManager {
	return &SessionManager{
		request: request,
		backend: backend,
		writer:  writer,
		db:      db,
	}
}

func (sm *SessionManager) VC() *ViewerContext {
	cookie, err := sm.request.Cookie(CookieName)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return EmptyVC()
		default:
			logrus.Errorf("Couldn't get cookie: %s", err.Error())
			return EmptyVC()
		}
	}

	token, err := ParseJWT(cookie.Value)
	if err != nil {
		logrus.Errorf("Couldn't parse JWT token string from cookie: %s", err.Error())
		return EmptyVC()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return EmptyVC()
	}

	sessionKeyClaim, ok := claims["session_key"]
	if !ok {
		return EmptyVC()
	}

	sessionKey, ok := sessionKeyClaim.(string)
	if !ok {
		return EmptyVC()
	}

	session, err := sm.backend.Get(sessionKey)
	if err != nil {
		switch err {
		case ErrNoActiveSession:
			return EmptyVC()
		default:
			logrus.Errorf("Couldn't call Get on backend object: %s", err.Error())
			return EmptyVC()
		}
	}

	user := &User{BaseModel: BaseModel{ID: session.UserID}}
	err = sm.db.Where(user).First(user).Error
	if err != nil {
		return EmptyVC()
	}

	return GenerateVC(user)
}

func (sm *SessionManager) MakeSessionForUserID(id uint) error {
	err := sm.backend.Create(id)
	if err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) MakeSessionForUser(u *User) error {
	return sm.MakeSessionForUserID(u.ID)
}

func (sm *SessionManager) Save() error {
	session, err := sm.backend.Session()
	if err != nil {
		return err
	}

	token, err := GenerateJWT(session.Key)
	if err != nil {
		return err
	}

	// Set proper flags on cookie for maximum security
	http.SetCookie(sm.writer, &http.Cookie{
		Name:  CookieName,
		Value: token,
	})

	return nil
}

func (sm *SessionManager) Destroy() error {
	if sm.backend != nil {
		err := sm.backend.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Session Backend API
////////////////////////////////////////////////////////////////////////////////

type SessionBackend interface {
	Get(key string) (*Session, error)
	Create(userID uint) error
	Destroy() error
	Session() (*Session, error)
	MarkAccessed() error
}

type BaseSessionBackend struct {
	session *Session
}

func (backend *BaseSessionBackend) Session() (*Session, error) {
	if backend.session == nil {
		return nil, ErrSessionNotCreated
	}
	return backend.session, nil
}

////////////////////////////////////////////////////////////////////////////////
// Session Backend Plugins
////////////////////////////////////////////////////////////////////////////////

type GormSessionBackend struct {
	BaseSessionBackend
	db *gorm.DB
}

func (s *GormSessionBackend) Get(key string) (*Session, error) {
	session := &Session{
		Key: key,
	}

	err := s.db.Where(session).First(session).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, ErrNoActiveSession
		default:
			return nil, err
		}
	}

	if time.Since(session.AccessedAt).Seconds() >= config.App.SessionExpirationSeconds {
		err = s.db.Delete(session).Error
		if err != nil {
			return nil, err
		}
		return nil, ErrSessionExpired
	}

	s.session = session

	err = s.MarkAccessed()
	if err != nil {
		return nil, err
	}

	return s.session, nil
}

func (s *GormSessionBackend) Create(userID uint) error {
	key := make([]byte, config.App.SessionKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	session := &Session{
		UserID: userID,
		Key:    base64.StdEncoding.EncodeToString(key),
	}

	err = s.db.Create(session).Error
	if err != nil {
		return err
	}
	s.session = session

	err = s.MarkAccessed()
	if err != nil {
		return err
	}

	return nil
}

func (s *GormSessionBackend) Destroy() error {
	if s.Session != nil {
		err := s.db.Delete(s.Session).Error
		if err != nil {
			return err
		}
		s.session = nil
	}

	return nil
}

func (s *GormSessionBackend) MarkAccessed() error {
	if s.session == nil {
		return ErrSessionNotCreated
	}
	s.session.AccessedAt = time.Now().UTC()
	return s.db.Save(s.session).Error
}
