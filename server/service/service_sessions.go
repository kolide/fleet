package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kolide/kolide/server/contexts/viewer"
	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/sso"
	"github.com/pkg/errors"
	"github.com/y0ssar1an/q"
)

func (svc service) InitiateSSO(ctx context.Context, idpID uint, relayURL, ssoHandle string) (string, error) {
	q.Q("initiate sso ")
	isProvider, err := svc.ds.IdentityProvider(idpID)
	if err != nil {
		return "", err
	}
	q.Q(isProvider)
	// get data about how to talk to the idp
	if isProvider.Metadata == "" {
		return "", errors.New("InitiateSSO missing metadata")
	}
	metadata, err := sso.ParseMetadata(isProvider.Metadata)
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO parsing metadata")
	}
	q.Q("got metadata")
	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting app config")
	}
	settings := sso.Settings{
		Metadata: metadata,
		// construct call back url to send to idp
		AssertionConsumerServiceURL: appConfig.KolideServerURL + "/api/v1/kolide/sso/callback",
	}
	idpURL, err := sso.CreateAuthorizationRequest(&settings, sso.RelayState(ssoHandle))
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO creating authorization")
	}
	q.Q("got url")
	// we're all ready to invoke idp with our redirect url, create a session in redis
	// so we can coordinate state across the idp, kolide, and the viewer, we set
	// session lifetime to five minutes,
	err = svc.ssoSessionStore.CreateSession(ssoHandle, relayURL, 300)
	if err != nil {
		return "", errors.Wrap(err, "creating sso session")
	}
	q.Q(idpURL)
	return idpURL, nil
}

func (svc service) CallbackSSO(ctx context.Context, ssoHandle, userID string) (string, error) {
	return "", nil
}

func (svc service) LoginSSO(ctx context.Context, ssoHandle string) (user *kolide.User, token string, err error) {
	return
}

func (svc service) Login(ctx context.Context, username, password string) (*kolide.User, string, error) {
	user, err := svc.userByEmailOrUsername(username)
	if _, ok := err.(kolide.NotFoundError); ok {
		return nil, "", authError{reason: "no such user"}
	}
	if err != nil {
		return nil, "", err
	}
	if !user.Enabled {
		return nil, "", authError{reason: "account disabled", clientReason: "account disabled"}
	}
	if err = user.ValidatePassword(password); err != nil {
		return nil, "", authError{reason: "bad password"}
	}
	token, err := svc.makeSession(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (svc service) userByEmailOrUsername(username string) (*kolide.User, error) {
	if strings.Contains(username, "@") {
		return svc.ds.UserByEmail(username)
	}
	return svc.ds.User(username)
}

// makeSession is a helper that creates a new session after authentication
func (svc service) makeSession(id uint) (string, error) {
	sessionKeySize := svc.config.Session.KeySize
	key := make([]byte, sessionKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	session := &kolide.Session{
		UserID:     id,
		Key:        base64.StdEncoding.EncodeToString(key),
		AccessedAt: time.Now().UTC(),
	}

	session, err = svc.ds.NewSession(session)
	if err != nil {
		return "", errors.Wrap(err, "creating new session")
	}

	tokenString, err := generateJWT(session.Key, svc.config.Auth.JwtKey)
	if err != nil {
		return "", errors.Wrap(err, "generating JWT token")
	}

	return tokenString, nil
}

func (svc service) Logout(ctx context.Context) error {
	// this should not return an error if the user wasn't logged in
	return svc.DestroySession(ctx)
}

func (svc service) DestroySession(ctx context.Context) error {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return errNoContext
	}

	session, err := svc.ds.SessionByID(vc.SessionID())
	if err != nil {
		return err
	}

	return svc.ds.DestroySession(session)
}

func (svc service) GetInfoAboutSessionsForUser(ctx context.Context, id uint) ([]*kolide.Session, error) {
	var validatedSessions []*kolide.Session

	sessions, err := svc.ds.ListSessionsForUser(id)
	if err != nil {
		return validatedSessions, err
	}

	for _, session := range sessions {
		if svc.validateSession(session) == nil {
			validatedSessions = append(validatedSessions, session)
		}
	}

	return validatedSessions, nil
}

func (svc service) DeleteSessionsForUser(ctx context.Context, id uint) error {
	return svc.ds.DestroyAllSessionsForUser(id)
}

func (svc service) GetInfoAboutSession(ctx context.Context, id uint) (*kolide.Session, error) {
	session, err := svc.ds.SessionByID(id)
	if err != nil {
		return nil, err
	}

	err = svc.validateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (svc service) GetSessionByKey(ctx context.Context, key string) (*kolide.Session, error) {
	session, err := svc.ds.SessionByKey(key)
	if err != nil {
		return nil, err
	}

	err = svc.validateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (svc service) DeleteSession(ctx context.Context, id uint) error {
	session, err := svc.ds.SessionByID(id)
	if err != nil {
		return err
	}
	return svc.ds.DestroySession(session)
}

func (svc service) validateSession(session *kolide.Session) error {
	if session == nil {
		return authError{
			reason:       "active session not present",
			clientReason: "session error",
		}
	}

	sessionDuration := svc.config.Session.Duration
	// duration 0 = unlimited
	if sessionDuration != 0 && time.Since(session.AccessedAt) >= sessionDuration {
		err := svc.ds.DestroySession(session)
		if err != nil {
			return errors.Wrap(err, "destroying session")
		}
		return authError{
			reason:       "expired session",
			clientReason: "session error",
		}
	}

	return svc.ds.MarkSessionAccessed(session)
}

// Given a session key create a JWT to be delivered to the client
func generateJWT(sessionKey, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_key": sessionKey,
	})

	return token.SignedString([]byte(jwtKey))
}
