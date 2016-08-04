package main

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

var (
	ErrNoActiveSession   = errors.New("Active session is not present in the database")
	ErrSessionNotCreated = errors.New("The session has not been created")
)

type ActiveSession struct {
	BaseModel
	UserID uint   `gorm:"not null"`
	Key    string `gorm:"not null;unique_index:idx_session_unique_key"`
}

////////////////////////////////////////////////////////////////////////////////
// Managing sessions
////////////////////////////////////////////////////////////////////////////////

type SessionManager struct {
	backend SessionBackend
	request *http.Request
	writer  http.ResponseWriter
	session *ActiveSession
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
	cookie, err := sm.request.Cookie("KolideSession")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return EmptyVC()
		default:
			logrus.Errorf("Couldn't get KolideSession cookie: %s", err.Error())
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

	token, err := GenerateJWTSession(session.Key)
	if err != nil {
		return err
	}

	http.SetCookie(sm.writer, &http.Cookie{
		Name:  "KolideSession",
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
	Get(key string) (*ActiveSession, error)
	Create(userID uint) error
	Destroy() error
	Session() (*ActiveSession, error)
}

type BaseSessionBackend struct {
	session *ActiveSession
}

func (backend *BaseSessionBackend) Session() (*ActiveSession, error) {
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

func (s *GormSessionBackend) Get(key string) (*ActiveSession, error) {
	session := &ActiveSession{
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
	s.session = session

	return s.session, nil
}

func (s *GormSessionBackend) Create(userID uint) error {
	session := &ActiveSession{
		UserID: userID,
		Key:    generateRandomText(32),
	}

	err := s.db.Preload("Users").Create(session).Error
	if err != nil {
		return err
	}
	s.session = session

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
