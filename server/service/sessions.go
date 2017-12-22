package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/sso"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) SSOSettings(ctx context.Context) (*kolide.SSOSettings, error) {
	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "SSOSettings getting app config")
	}
	settings := &kolide.SSOSettings{
		IDPName:     appConfig.IDPName,
		IDPImageURL: appConfig.IDPImageURL,
		SSOEnabled:  appConfig.EnableSSO,
	}
	return settings, nil
}

func (svc service) InitiateSSO(ctx context.Context, redirectURL string) (string, error) {
	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting app config")
	}

	metadata, err := svc.getMetadata(appConfig)
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting metadata")
	}

	settings := sso.Settings{
		Metadata: metadata,
		// Construct call back url to send to idp
		AssertionConsumerServiceURL: appConfig.KolideServerURL + "/api/v1/kolide/sso/callback",
		SessionStore:                svc.ssoSessionStore,
		OriginalURL:                 redirectURL,
	}

	// If issuer is not explicitly set, default to host name.
	var issuer string
	if appConfig.EntityID == "" {
		u, err := url.Parse(appConfig.KolideServerURL)
		if err != nil {
			return "", errors.Wrap(err, "parsing kolide server url")
		}
		issuer = u.Hostname()
	} else {
		issuer = appConfig.EntityID
	}
	idpURL, err := sso.CreateAuthorizationRequest(&settings, issuer)
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO creating authorization")
	}

	return idpURL, nil
}

func (svc service) getMetadata(config *kolide.AppConfig) (*sso.Metadata, error) {
	if config.MetadataURL != "" {
		metadata, err := sso.GetMetadata(config.MetadataURL, svc.metaDataClient)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}
	if config.Metadata != "" {
		metadata, err := sso.ParseMetadata(config.Metadata)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}
	return nil, errors.Errorf("missing metadata for idp %s", config.IDPName)
}

func (svc service) CallbackSSO(ctx context.Context, auth kolide.Auth) (*kolide.SSOSession, error) {
	// The signature and validity of auth response has been checked already in
	// validation middleware.
	sess, err := svc.ssoSessionStore.Get(auth.RequestID())
	if err != nil {
		return nil, errors.Wrap(err, "fetching sso session in callback")
	}
	// Remove session to so that is can't be reused before it expires.
	err = svc.ssoSessionStore.Expire(auth.RequestID())
	if err != nil {
		return nil, errors.Wrap(err, "expiring sso session in callback")
	}
	user, err := svc.userByEmailOrUsername(auth.UserID())
	if err != nil {
		return nil, errors.Wrap(err, "finding user in sso callback")
	}
	// if user is not active they are not authorized to use the application
	if !user.Enabled || user.Deleted {
		return nil, errors.New("user authorization failed")
	}
	// if the user is not sso enabled they are not authorized
	if !user.SSOEnabled {
		return nil, errors.New("user not configured to use sso")
	}
	token, err := svc.makeSession(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "making user session in sso callback")
	}
	result := &kolide.SSOSession{
		Token:       token,
		RedirectURL: sess.OriginalURL,
	}
	if !strings.HasPrefix(result.RedirectURL, "/") {
		result.RedirectURL = "/" + result.RedirectURL
	}
	return result, nil
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
	if user.SSOEnabled {
		const errMessage = "password login not allowed for single sign on users"
		return nil, "", authError{reason: errMessage, clientReason: errMessage}
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

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeGetInfoAboutSessionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getInfoAboutSessionRequest{ID: id}, nil
}

func decodeGetInfoAboutSessionsForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getInfoAboutSessionsForUserRequest{ID: id}, nil
}

func decodeDeleteSessionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteSessionRequest{ID: id}, nil
}

func decodeDeleteSessionsForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteSessionsForUserRequest{ID: id}, nil
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Username = strings.ToLower(req.Username)
	return req, nil
}

func decodeInitiateSSORequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req initiateSSORequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCallbackSSORequest(ctx context.Context, r *http.Request) (interface{}, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, errors.Wrap(err, "decode sso callback")
	}
	authResponse, err := sso.DecodeAuthResponse(r.FormValue("SAMLResponse"))
	if err != nil {
		return nil, errors.Wrap(err, "decoding sso callback")
	}
	return authResponse, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type loginRequest struct {
	Username string // can be username or email
	Password string
}

type loginResponse struct {
	User  *kolide.User `json:"user,omitempty"`
	Token string       `json:"token,omitempty"`
	Err   error        `json:"error,omitempty"`
}

func (r loginResponse) error() error { return r.Err }

func makeLoginEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		user, token, err := svc.Login(ctx, req.Username, req.Password)
		if err != nil {
			return loginResponse{Err: err}, nil
		}
		return loginResponse{user, token, nil}, nil
	}
}

type logoutResponse struct {
	Err error `json:"error,omitempty"`
}

func (r logoutResponse) error() error { return r.Err }

func makeLogoutEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := svc.Logout(ctx)
		if err != nil {
			return logoutResponse{Err: err}, nil
		}
		return logoutResponse{}, nil
	}
}

type getInfoAboutSessionRequest struct {
	ID uint
}

type getInfoAboutSessionResponse struct {
	SessionID uint      `json:"session_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Err       error     `json:"error,omitempty"`
}

func (r getInfoAboutSessionResponse) error() error { return r.Err }

func makeGetInfoAboutSessionEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getInfoAboutSessionRequest)
		session, err := svc.GetInfoAboutSession(ctx, req.ID)
		if err != nil {
			return getInfoAboutSessionResponse{Err: err}, nil
		}

		return getInfoAboutSessionResponse{
			SessionID: session.ID,
			UserID:    session.UserID,
			CreatedAt: session.CreatedAt,
		}, nil
	}
}

type getInfoAboutSessionsForUserRequest struct {
	ID uint
}

type getInfoAboutSessionsForUserResponse struct {
	Sessions []getInfoAboutSessionResponse `json:"sessions"`
	Err      error                         `json:"error,omitempty"`
}

func (r getInfoAboutSessionsForUserResponse) error() error { return r.Err }

func makeGetInfoAboutSessionsForUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getInfoAboutSessionsForUserRequest)
		sessions, err := svc.GetInfoAboutSessionsForUser(ctx, req.ID)
		if err != nil {
			return getInfoAboutSessionsForUserResponse{Err: err}, nil
		}
		var resp getInfoAboutSessionsForUserResponse
		for _, session := range sessions {
			resp.Sessions = append(resp.Sessions, getInfoAboutSessionResponse{
				SessionID: session.ID,
				UserID:    session.UserID,
				CreatedAt: session.CreatedAt,
			})
		}
		return resp, nil
	}
}

type deleteSessionRequest struct {
	ID uint
}

type deleteSessionResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteSessionResponse) error() error { return r.Err }

func makeDeleteSessionEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteSessionRequest)
		err := svc.DeleteSession(ctx, req.ID)
		if err != nil {
			return deleteSessionResponse{Err: err}, nil
		}
		return deleteSessionResponse{}, nil
	}
}

type deleteSessionsForUserRequest struct {
	ID uint
}

type deleteSessionsForUserResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteSessionsForUserResponse) error() error { return r.Err }

func makeDeleteSessionsForUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteSessionsForUserRequest)
		err := svc.DeleteSessionsForUser(ctx, req.ID)
		if err != nil {
			return deleteSessionsForUserResponse{Err: err}, nil
		}
		return deleteSessionsForUserResponse{}, nil
	}
}

type initiateSSORequest struct {
	RelayURL string `json:"relay_url"`
}

type initiateSSOResponse struct {
	URL string `json:"url,omitempty"`
	Err error  `json:"error,omitempty"`
}

func (r initiateSSOResponse) error() error { return r.Err }

func makeInitiateSSOEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(initiateSSORequest)
		idProviderURL, err := svc.InitiateSSO(ctx, req.RelayURL)
		if err != nil {
			return initiateSSOResponse{Err: err}, nil
		}
		return initiateSSOResponse{URL: idProviderURL}, nil
	}
}

type callbackSSOResponse struct {
	content string
	Err     error `json:"error,omitempty"`
}

func (r callbackSSOResponse) error() error { return r.Err }

// If html is present we return a web page
func (r callbackSSOResponse) html() string { return r.content }

func makeCallbackSSOEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		authResponse := request.(kolide.Auth)
		session, err := svc.CallbackSSO(ctx, authResponse)
		var resp callbackSSOResponse
		if err != nil {
			// redirect to login page on front end if there was some problem,
			// errors should still be logged
			session = &kolide.SSOSession{
				RedirectURL: "/login",
				Token:       "",
			}
			resp.Err = err
		}
		relayStateLoadPage := ` <html>
     <script type='text/javascript'>
     var redirectURL = {{.RedirectURL}};
     window.localStorage.setItem('KOLIDE::auth_token', '{{.Token}}');
     window.location = redirectURL;
     </script>
     <body>
     Redirecting to Kolide...
     </body>
     </html>
    `
		tmpl, err := template.New("relayStateLoader").Parse(relayStateLoadPage)
		if err != nil {
			return nil, err
		}
		var writer bytes.Buffer
		err = tmpl.Execute(&writer, session)
		if err != nil {
			return nil, err
		}
		resp.content = writer.String()
		return resp, nil
	}
}

type ssoSettingsResponse struct {
	Settings *kolide.SSOSettings `json:"settings,omitempty"`
	Err      error               `json:"error,omitempty"`
}

func (r ssoSettingsResponse) error() error { return r.Err }

func makeSSOSettingsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, unused interface{}) (interface{}, error) {
		settings, err := svc.SSOSettings(ctx)
		if err != nil {
			return ssoSettingsResponse{Err: err}, nil
		}
		return ssoSettingsResponse{Settings: settings}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) SSOSettings(ctx context.Context) (settings *kolide.SSOSettings, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SSOSettings", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	settings, err = mw.Service.SSOSettings(ctx)
	return
}

func (mw metricsMiddleware) InitiateSSO(ctx context.Context, relayValue string) (idpURL string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "InitiateSSO", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	idpURL, err = mw.Service.InitiateSSO(ctx, relayValue)
	return
}

func (mw metricsMiddleware) CallbackSSO(ctx context.Context, auth kolide.Auth) (sess *kolide.SSOSession, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CallbackSSO", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	sess, err = mw.Service.CallbackSSO(ctx, auth)
	return
}

func (mw metricsMiddleware) Login(ctx context.Context, username string, password string) (*kolide.User, string, error) {
	var (
		user  *kolide.User
		token string
		err   error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	user, token, err = mw.Service.Login(ctx, username, password)
	return user, token, err
}

func (mw metricsMiddleware) Logout(ctx context.Context) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "Logout", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.Logout(ctx)
	return err
}

func (mw metricsMiddleware) DestroySession(ctx context.Context) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "DestroySession", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.DestroySession(ctx)
	return err
}

func (mw metricsMiddleware) GetInfoAboutSessionsForUser(ctx context.Context, id uint) ([]*kolide.Session, error) {
	var (
		sessions []*kolide.Session
		err      error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "GetInfoAboutSessionsForUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	sessions, err = mw.Service.GetInfoAboutSessionsForUser(ctx, id)
	return sessions, err
}

func (mw metricsMiddleware) DeleteSessionsForUser(ctx context.Context, id uint) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteSessionsForUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.DeleteSessionsForUser(ctx, id)
	return err
}

func (mw metricsMiddleware) GetInfoAboutSession(ctx context.Context, id uint) (*kolide.Session, error) {
	var (
		session *kolide.Session
		err     error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "GetInfoAboutSession", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	session, err = mw.Service.GetInfoAboutSession(ctx, id)
	return session, err
}

func (mw metricsMiddleware) GetSessionByKey(ctx context.Context, key string) (*kolide.Session, error) {
	var (
		session *kolide.Session
		err     error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "GetSessionByKey", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	session, err = mw.Service.GetSessionByKey(ctx, key)
	return session, err
}

func (mw metricsMiddleware) DeleteSession(ctx context.Context, id uint) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteSession", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.DeleteSession(ctx, id)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) Login(ctx context.Context, username, password string) (user *kolide.User, token string, err error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Login",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, token, err = mw.Service.Login(ctx, username, password)
	return
}

func (mw loggingMiddleware) Logout(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Logout",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.Logout(ctx)
	return
}

func (mw loggingMiddleware) InitiateSSO(ctx context.Context, relayURL string) (idpURL string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "InitiateSSO",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	idpURL, err = mw.Service.InitiateSSO(ctx, relayURL)
	return
}

func (mw loggingMiddleware) CallbackSSO(ctx context.Context, auth kolide.Auth) (sess *kolide.SSOSession, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "CallbackSSO",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	sess, err = mw.Service.CallbackSSO(ctx, auth)
	return
}

func (mw loggingMiddleware) SSOSettings(ctx context.Context) (settings *kolide.SSOSettings, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SSOSettings",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	settings, err = mw.Service.SSOSettings(ctx)
	return
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func (mw validationMiddleware) CallbackSSO(ctx context.Context, auth kolide.Auth) (*kolide.SSOSession, error) {
	invalid := &invalidArgumentError{}
	session, err := mw.ssoSessionStore.Get(auth.RequestID())
	if err != nil {
		invalid.Append("session", "missing for request")
		return nil, invalid
	}
	validator, err := sso.NewValidator(session.Metadata)
	if err != nil {
		return nil, errors.Wrap(err, "creating validator from metadata")
	}
	// make sure the response hasn't been tampered with
	auth, err = validator.ValidateSignature(auth)
	if err != nil {
		invalid.Appendf("sso response", "signature validation failed %s", err.Error())
		return nil, invalid
	}
	// make sure the response isn't stale
	err = validator.ValidateResponse(auth)
	if err != nil {
		invalid.Appendf("sso response", "response validation failed %s", err.Error())
	}

	return mw.Service.CallbackSSO(ctx, auth)
}
