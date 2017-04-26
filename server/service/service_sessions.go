package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kolide/kolide/server/contexts/viewer"
	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/sso"
	"github.com/pkg/errors"
)

const ssoSessionExpiry = 300 // five minutes to complete sso login process

func (svc service) InitiateSSO(ctx context.Context, idpID uint, relayURL, ssoHandle string) (string, error) {
	isProvider, err := svc.ds.IdentityProvider(idpID)
	if err != nil {
		return "", err
	}
	metadata, err := getMetadata(isProvider)

	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting metadata")
	}
	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting app config")
	}
	settings := sso.Settings{
		Metadata: metadata,
		// Construct call back url to send to idp
		AssertionConsumerServiceURL: appConfig.KolideServerURL + "/api/v1/kolide/sso/callback",
	}
	// encrypt sso handle before sending it to IDP
	encryptedSSOHandle, err := svc.ssoSessionStore.EncryptSSOHandle(ssoHandle, appConfig.AESKey, ssoSessionExpiry)
	if err != nil {
		return "", errors.Wrap(err, "encrypting sso session handle")
	}
	// If issuer is not explicitly set, default to host name.
	var issuer string
	if isProvider.IssuerURI == "" {
		u, err := url.Parse(appConfig.KolideServerURL)
		if err != nil {
			return "", errors.Wrap(err, "parsing kolide server url")
		}
		issuer = u.Hostname()
	} else {
		issuer = isProvider.IssuerURI
	}
	idpURL, err := sso.CreateAuthorizationRequest(&settings, issuer, sso.RelayState(encryptedSSOHandle))
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO creating authorization")
	}
	// We're all ready to invoke idp with our redirect url, create a session in redis
	// so we can coordinate state across the idp, kolide, and the viewer, we set
	// session lifetime to five minutes,
	err = svc.ssoSessionStore.CreateSession(ssoHandle, relayURL, ssoSessionExpiry)
	if err != nil {
		return "", errors.Wrap(err, "creating sso session")
	}
	return idpURL, nil
}

func getMetadata(idp *kolide.IdentityProvider) (*sso.IDPMetadata, error) {
	if idp.MetadataURL != "" {
		metadata, err := sso.GetMetadata(idp.MetadataURL, 5*time.Second)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}
	if idp.Metadata != "" {
		metadata, err := sso.ParseMetadata(idp.Metadata)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}
	return nil, errors.Errorf("missing metadata for idp %s", idp.Name)
}

func (svc service) CallbackSSO(ctx context.Context, auth kolide.Auth) (string, error) {
	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return "", errors.Wrap(err, "sso authentication callback")
	}
	status, err := auth.Status()
	if err != nil {
		return "", errors.Wrap(err, "fetching sso response status")
	}
	if status != sso.Success {
		svc.logger.Log(
			"method", "CallbackSSO",
			"err", auth.StatusDescription(),
		)
		// TODO: Create custom 401 page
		return appConfig.KolideServerURL + "/404", nil
	}
	ssoHandle, err := svc.ssoSessionStore.DecryptSSOHandle(auth.RelayState(), appConfig.AESKey)
	if err != nil {
		return "", errors.Wrap(err, "deciphering sso handle")
	}
	// Setting user id that we get from the IDP indicates we have successfully
	// authenticated.
	sess, err := svc.ssoSessionStore.UpdateSession(ssoHandle, auth.UserID())
	// Because we're being called from an external IDP (via the browser) as opposed to returning json from
	// an api call, the only thing that probably could go wrong is that the session expired
	// because the user sat with the login page too long.  In this case, we'll reload
	// the original resource in the SPA, which will call LoginSSO, which will fail because
	// the session expired and we can display a more useful message.
	svc.logger.Log("method", "CallbackSSO", "err", err)
	if strings.HasPrefix(sess.OriginalURL, "/") {
		return appConfig.KolideServerURL + sess.OriginalURL, nil
	}
	return appConfig.KolideServerURL + "/" + sess.OriginalURL, nil
}

func (svc service) LoginSSO(ctx context.Context, ssoHandle string) (user *kolide.User, token string, err error) {
	var session *sso.Session
	session, err = svc.ssoSessionStore.GetSession(ssoHandle)
	if err != nil {
		return nil, "", errors.Wrap(err, "sso login retrieving session")
	}
	defer func() { err = svc.ssoSessionStore.ExpireSession(ssoHandle) }()
	user, err = svc.userByEmailOrUsername(session.UserName)
	if err != nil {
		return nil, "", errors.Wrap(err, "fetching username in sso login")
	}
	token, err = svc.makeSession(user.ID)
	return user, token, nil
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
