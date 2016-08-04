package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	// An error returned by SessionBackend.Get() if no session record was found
	// in the database
	ErrNoActiveSession = errors.New("Active session is not present in the database")

	// An error returned by SessionBackend methods when no session object has
	// been created yet but the requested action requires one
	ErrSessionNotCreated = errors.New("The session has not been created")

	// An error returned by SessionBackend.Get() when a session is requested but
	// it has expired
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

// SessionManager is a management object which helps with the administration of
// sessions within the application. Use NewSessionManager to create an instance
type SessionManager struct {
	backend SessionBackend
	request *http.Request
	writer  http.ResponseWriter
	session *Session
	vc      *ViewerContext
	db      *gorm.DB
}

// NewSessionManager allows you to get a SessionManager instance for a given
// web request. Unless you're interacting with login, logout, or core auth
// code, this should be abstracted by the ViewerContext pattern.
func NewSessionManager(request *http.Request, writer http.ResponseWriter, backend SessionBackend, db *gorm.DB) *SessionManager {
	return &SessionManager{
		request: request,
		backend: backend,
		writer:  writer,
		db:      db,
	}
}

// Get the ViewerContext instance for a user represented by the active session
func (sm *SessionManager) VC() *ViewerContext {
	if sm.session == nil {
		cookie, err := sm.request.Cookie(CookieName)
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				// No cookie was set
				return EmptyVC()
			default:
				// Something went wrong and the cookie may or may not be set
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
			logrus.Error("Could not parse the claims from the JWT token")
			return EmptyVC()
		}

		sessionKeyClaim, ok := claims["session_key"]
		if !ok {
			logrus.Warn("JWT did not have session_key claim")
			return EmptyVC()
		}

		sessionKey, ok := sessionKeyClaim.(string)
		if !ok {
			logrus.Warn("JWT session_key claim was not a string")
			return EmptyVC()
		}

		session, err := sm.backend.Get(sessionKey)
		if err != nil {
			switch err {
			case ErrNoActiveSession:
				// If the code path got this far, it's likely that the user was logged
				// in some time in the past, but their session has been expired since
				// their last usage of the application
				return EmptyVC()
			default:
				logrus.Errorf("Couldn't call Get on backend object: %s", err.Error())
				return EmptyVC()
			}
		}
		sm.session = session
	}

	if sm.vc == nil {
		// Generating a VC requires a user struct. Attempt to populate one using
		// the user id of the current session holder
		user := &User{BaseModel: BaseModel{ID: sm.session.UserID}}
		err := sm.db.Where(user).First(user).Error
		if err != nil {
			return EmptyVC()
		}

		sm.vc = GenerateVC(user)
	}

	return sm.vc
}

// MakeSessionForUserID creates a session in the database for a given user id.
// You must call Save() after calling this.
func (sm *SessionManager) MakeSessionForUserID(id uint) error {
	err := sm.backend.Create(id)
	if err != nil {
		return err
	}
	return nil
}

// MakeSessionForUserID creates a session in the database for a given user
// You must call Save() after calling this.
func (sm *SessionManager) MakeSessionForUser(u *User) error {
	return sm.MakeSessionForUserID(u.ID)
}

// Save writes the current session to a token and delivers the token as a cookie
// to the user. Save must be called after every write action on this struct
// (MakeSessionForUser, Destroy, etc.)
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

// Destory deletes the active session from the database and erases the session
// instance from this object's access. You must call Save() after calling this.
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

// SessionBackend is the abstract interface that all session backends must
// conform to. SessionBackend instances are only expected to exist within the
// context of a single request.
type SessionBackend interface {
	// Given a session key, find and return a session object or an error if one
	// could not be found for the given key
	Get(key string) (*Session, error)

	// Create a session object tied to the given user ID
	Create(userID uint) error

	// Destroy the currently tracked session
	Destroy() error

	// Return the currently tracked session
	Session() (*Session, error)

	// Mark the currently tracked session as access to extend expiration
	MarkAccessed() error
}

// BaseSessionBackend is a convenience struct that all SessionBackends can
// embed
type BaseSessionBackend struct {
	session *Session
}

// Session returns the currently active session or an error if one is not
// available
func (backend *BaseSessionBackend) Session() (*Session, error) {
	if backend.session == nil {
		return nil, ErrSessionNotCreated
	}
	return backend.session, nil
}

////////////////////////////////////////////////////////////////////////////////
// Session Backend Plugins
////////////////////////////////////////////////////////////////////////////////

// GormSessionBackend stores sessions using a pre-instantiated gorm database
// object
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
	if _, err := s.Session(); err != nil {
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

////////////////////////////////////////////////////////////////////////////////
// Session management HTTP endpoints
////////////////////////////////////////////////////////////////////////////////

type DeleteSessionRequestBody struct {
	SessionID uint `json:"session_id" binding:"required"`
}

func DeleteSession(c *gin.Context) {
	var body DeleteSessionRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	db, err := GetDB(c)
	if err != nil {
		logrus.Errorf("Could not open database: %s", err.Error())
		DatabaseError(c)
		return
	}

	vc, err := VC(c, db)
	if err != nil {
		logrus.Errorf("Could not create VC: %s", err.Error())
		DatabaseError(c) // TODO tampered?
		return
	}

	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	session := &Session{
		BaseModel: BaseModel{
			ID: body.SessionID,
		},
	}
	err = db.Where(session).First(session).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(404, nil)
			return
		default:
			DatabaseError(c)
			return
		}
	}

	user := &User{
		BaseModel: BaseModel{
			ID: session.UserID,
		},
	}
	err = db.Where(user).First(user).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	err = db.Delete(session).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	c.JSON(200, nil)
}

type DeleteSessionsForUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func DeleteSessionsForUser(c *gin.Context) {
	var body DeleteSessionsForUserRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	db, err := GetDB(c)
	if err != nil {
		logrus.Errorf("Could not open database: %s", err.Error())
		DatabaseError(c)
		return
	}

	vc, err := VC(c, db)
	if err != nil {
		logrus.Errorf("Could not create VC: %s", err.Error())
		DatabaseError(c) // TODO tampered?
		return
	}

	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	var user User
	user.ID = body.ID
	user.Username = body.Username
	err = db.Where(&user).First(&user).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	if !vc.CanPerformWriteActionOnUser(&user) {
		UnauthorizedError(c)
		return
	}

	err = db.Delete(&Session{}, "user_id = ?", user.ID).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	c.JSON(200, nil)

}

type GetInfoAboutSessionRequestBody struct {
	SessionKey string `json:"session_key" binding:"required"`
}

type SessionInfoResponseBody struct {
	SessionID  uint      `json:"session_id"`
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	AccessedAt time.Time `json:"created_at"`
}

func GetInfoAboutSession(c *gin.Context) {
	var body GetInfoAboutSessionRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	db, err := GetDB(c)
	if err != nil {
		logrus.Errorf("Could not open database: %s", err.Error())
		DatabaseError(c)
		return
	}

	vc, err := VC(c, db)
	if err != nil {
		logrus.Errorf("Could not create VC: %s", err.Error())
		DatabaseError(c) // TODO tampered?
		return
	}

	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	var session Session
	session.Key = body.SessionKey
	err = db.Where(&session).First(&session).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	var user User
	user.ID = session.UserID
	err = db.Where(&user).First(&user).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	if !vc.IsAdmin() && !vc.IsUserID(user.ID) {
		UnauthorizedError(c)
		return
	}

	c.JSON(200, &SessionInfoResponseBody{
		SessionID:  session.ID,
		UserID:     session.UserID,
		CreatedAt:  session.CreatedAt,
		AccessedAt: session.AccessedAt,
	})
}

type GetInfoAboutSessionsForUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type GetInfoAboutSessionsForUserResponseBody struct {
	Sessions []*SessionInfoResponseBody `json:"sessions"`
}

func GetInfoAboutSessionsForUser(c *gin.Context) {
	var body GetInfoAboutSessionsForUserRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	db, err := GetDB(c)
	if err != nil {
		logrus.Errorf("Could not open database: %s", err.Error())
		DatabaseError(c)
		return
	}

	vc, err := VC(c, db)
	if err != nil {
		logrus.Errorf("Could not create VC: %s", err.Error())
		DatabaseError(c) // TODO tampered?
		return
	}

	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	var user User
	user.ID = body.ID
	user.Username = body.Username
	err = db.Where(&user).First(&user).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	if !vc.IsAdmin() && !vc.IsUserID(user.ID) {
		UnauthorizedError(c)
		return
	}

	var sessions []*Session
	err = db.Where("user_id = ?", user.ID).Find(&sessions).Error
	if err != nil {
		DatabaseError(c)
		return
	}

	var response []*SessionInfoResponseBody
	for _, session := range sessions {
		response = append(response, &SessionInfoResponseBody{
			SessionID:  session.ID,
			UserID:     session.UserID,
			CreatedAt:  session.CreatedAt,
			AccessedAt: session.AccessedAt,
		})
	}

	c.JSON(200, &GetInfoAboutSessionsForUserResponseBody{
		Sessions: response,
	})
}
