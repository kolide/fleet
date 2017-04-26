package kolide

import (
	"context"
	"time"
)

// SessionStore is the abstract interface that all session backends must
// conform to.
type SessionStore interface {
	// Given a session key, find and return a session object or an error if one
	// could not be found for the given key
	SessionByKey(key string) (*Session, error)

	// Given a session id, find and return a session object or an error if one
	// could not be found for the given id
	SessionByID(id uint) (*Session, error)

	// Find all of the active sessions for a given user
	ListSessionsForUser(id uint) ([]*Session, error)

	// Store a new session struct
	NewSession(session *Session) (*Session, error)

	// Destroy the currently tracked session
	DestroySession(session *Session) error

	// Destroy all of the sessions for a given user
	DestroyAllSessionsForUser(id uint) error

	// Mark the currently tracked session as access to extend expiration
	MarkSessionAccessed(session *Session) error
}

type Auth interface {
	RelayState() string
	UserID() string
	Status() (int, error)
	StatusDescription() string
}

type SessionService interface {
	// InitiateSSO is used to initiate an SSO session and returns a URL that
	// can be used in a redirect to the IDP.
	// Arguments: idpID is the database id of the Identity Provider, relayValue is a token
	// that is stored in the browser state and is associated with the url that the
	// user was accessing when prompted to log in.  The ssoHandle is a unique value generated in the front end
	// and is used to reference information from the SSO process we initiate.
	InitiateSSO(ctx context.Context, idpID uint, relayValue, ssoHandle string) (string, error)
	// CallbackSSO handles the IDP response.  The original URL the viewer attempted
	// to access is returned from this function so we can redirect back to the front end and
	// load the page the viewer originally attempted to access when prompted for login.
	CallbackSSO(ctx context.Context, auth Auth) (string, error)
	// LoginSSO is invoked from the front end after a successful SSO/SAML transaction
	// to retrieve session information.
	LoginSSO(ctx context.Context, ssoHandle string) (user *User, token string, err error)
	Login(ctx context.Context, username, password string) (user *User, token string, err error)
	Logout(ctx context.Context) (err error)
	DestroySession(ctx context.Context) (err error)
	GetInfoAboutSessionsForUser(ctx context.Context, id uint) (sessions []*Session, err error)
	DeleteSessionsForUser(ctx context.Context, id uint) (err error)
	GetInfoAboutSession(ctx context.Context, id uint) (session *Session, err error)
	GetSessionByKey(ctx context.Context, key string) (session *Session, err error)
	DeleteSession(ctx context.Context, id uint) (err error)
}

// Session is the model object which represents what an active session is
type Session struct {
	CreateTimestamp
	ID         uint
	AccessedAt time.Time `db:"accessed_at"`
	UserID     uint      `db:"user_id"`
	Key        string
}
