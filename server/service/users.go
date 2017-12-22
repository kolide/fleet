package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	invite, err := svc.VerifyInvite(ctx, *p.InviteToken)
	if err != nil {
		return nil, err
	}

	// set the payload Admin property based on an existing invite.
	p.Admin = &invite.Admin

	user, err := svc.newUser(p)
	if err != nil {
		return nil, err
	}

	err = svc.ds.DeleteInvite(invite.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc service) NewAdminCreatedUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	return svc.newUser(p)
}

func (svc service) newUser(p kolide.UserPayload) (*kolide.User, error) {
	var ssoEnabled bool
	// if user is SSO generate a fake password
	if p.SSOInvite != nil && *p.SSOInvite {
		fakePassword, err := generateRandomText(14)
		if err != nil {
			return nil, err
		}
		p.Password = &fakePassword
		ssoEnabled = true
	}
	user, err := p.User(svc.config.Auth.SaltKeySize, svc.config.Auth.BcryptCost)
	if err != nil {
		return nil, err
	}
	user.SSOEnabled = ssoEnabled
	user, err = svc.ds.NewUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc service) ChangeUserAdmin(ctx context.Context, id uint, isAdmin bool) (*kolide.User, error) {
	user, err := svc.ds.UserByID(id)
	if err != nil {
		return nil, err
	}
	user.Admin = isAdmin
	if err = svc.saveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (svc service) ChangeUserEnabled(ctx context.Context, id uint, isEnabled bool) (*kolide.User, error) {
	user, err := svc.ds.UserByID(id)
	if err != nil {
		return nil, err
	}
	user.Enabled = isEnabled
	if err = svc.saveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (svc service) ModifyUser(ctx context.Context, userID uint, p kolide.UserPayload) (*kolide.User, error) {
	user, err := svc.User(ctx, userID)
	if err != nil {
		return nil, err
	}

	// the method assumes that the correct authorization
	// has been validated higher up the stack
	if p.Admin != nil {
		user.Admin = *p.Admin
	}

	if p.Enabled != nil {
		user.Enabled = *p.Enabled
	}

	if p.Username != nil {
		user.Username = *p.Username
	}

	if p.Name != nil {
		user.Name = *p.Name
	}

	if p.Email != nil {
		err = svc.modifyEmailAddress(ctx, user, *p.Email, p.Password)
		if err != nil {
			return nil, err
		}
	}

	if p.Position != nil {
		user.Position = *p.Position
	}

	if p.GravatarURL != nil {
		user.GravatarURL = *p.GravatarURL
	}

	if p.SSOEnabled != nil {
		user.SSOEnabled = *p.SSOEnabled
	}

	err = svc.saveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc service) modifyEmailAddress(ctx context.Context, user *kolide.User, email string, password *string) error {
	// password requirement handled in validation middleware
	if password != nil {
		err := user.ValidatePassword(*password)
		if err != nil {
			return newPermissionError("password", "incorrect password")
		}
	}
	random, err := kolide.RandomText(svc.config.App.TokenKeySize)
	if err != nil {
		return err
	}
	token := base64.URLEncoding.EncodeToString([]byte(random))
	err = svc.ds.PendingEmailChange(user.ID, email, token)
	if err != nil {
		return err
	}
	config, err := svc.AppConfig(ctx)
	if err != nil {
		return err
	}
	changeEmail := kolide.Email{
		Subject: "Confirm Kolide Email Change",
		To:      []string{email},
		Config:  config,
		Mailer: &kolide.ChangeEmailMailer{
			Token:           token,
			KolideServerURL: template.URL(config.KolideServerURL),
		},
	}
	return svc.mailService.SendEmail(changeEmail)
}

func (svc service) ChangeUserEmail(ctx context.Context, token string) (string, error) {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return "", errNoContext
	}
	return svc.ds.ConfirmPendingEmailChange(vc.UserID(), token)
}

func (svc service) User(ctx context.Context, id uint) (*kolide.User, error) {
	return svc.ds.UserByID(id)
}

func (svc service) AuthenticatedUser(ctx context.Context) (*kolide.User, error) {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}
	if !vc.IsLoggedIn() {
		return nil, permissionError{}
	}
	return vc.User, nil
}

func (svc service) ListUsers(ctx context.Context, opt kolide.ListOptions) ([]*kolide.User, error) {
	return svc.ds.ListUsers(opt)
}

// setNewPassword is a helper for changing a user's password. It should be
// called to set the new password after proper authorization has been
// performed.
func (svc service) setNewPassword(ctx context.Context, user *kolide.User, password string) error {
	err := user.SetPassword(password, svc.config.Auth.SaltKeySize, svc.config.Auth.BcryptCost)
	if err != nil {
		return errors.Wrap(err, "setting new password")
	}
	if user.SSOEnabled {
		return errors.New("set password for single sign on user not allowed")
	}
	err = svc.saveUser(user)
	if err != nil {
		return errors.Wrap(err, "saving changed password")
	}

	return nil
}

func (svc service) ChangePassword(ctx context.Context, oldPass, newPass string) error {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return errNoContext
	}
	if vc.User.SSOEnabled {
		return errors.New("change password for single sign on user not allowed")
	}
	if err := vc.User.ValidatePassword(newPass); err == nil {
		return newInvalidArgumentError("new_password", "cannot reuse old password")
	}

	if err := vc.User.ValidatePassword(oldPass); err != nil {
		return newInvalidArgumentError("old_password", "old password does not match")
	}

	if err := svc.setNewPassword(ctx, vc.User, newPass); err != nil {
		return errors.Wrap(err, "setting new password")
	}
	return nil
}

func (svc service) ResetPassword(ctx context.Context, token, password string) error {
	reset, err := svc.ds.FindPassswordResetByToken(token)
	if err != nil {
		return errors.Wrap(err, "looking up reset by token")
	}
	user, err := svc.User(ctx, reset.UserID)
	if err != nil {
		return errors.Wrap(err, "retrieving user")
	}
	if user.SSOEnabled {
		return errors.New("password reset for single sign on user not allowed")
	}

	// prevent setting the same password
	if err := user.ValidatePassword(password); err == nil {
		return newInvalidArgumentError("new_password", "cannot reuse old password")
	}

	err = svc.setNewPassword(ctx, user, password)
	if err != nil {
		return errors.Wrap(err, "setting new password")
	}

	// delete password reset tokens for user
	if err := svc.ds.DeletePasswordResetRequestsForUser(user.ID); err != nil {
		return errors.Wrap(err, "deleting password reset requests")
	}

	// Clear sessions so that any other browsers will have to log in with
	// the new password
	if err := svc.DeleteSessionsForUser(ctx, user.ID); err != nil {
		return errors.Wrap(err, "deleting user sessions")
	}

	return nil
}

func (svc service) PerformRequiredPasswordReset(ctx context.Context, password string) (*kolide.User, error) {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}
	user := vc.User
	if user.SSOEnabled {
		return nil, errors.New("password reset for single sign on user not allowed")
	}
	if !user.AdminForcedPasswordReset {
		return nil, errors.New("user does not require password reset")
	}

	// prevent setting the same password
	if err := user.ValidatePassword(password); err == nil {
		return nil, newInvalidArgumentError("new_password", "cannot reuse old password")
	}

	user.AdminForcedPasswordReset = false
	err := svc.setNewPassword(ctx, user, password)
	if err != nil {
		return nil, errors.Wrap(err, "setting new password")
	}

	// Sessions should already have been cleared when the reset was
	// required

	return user, nil
}

func (svc service) RequirePasswordReset(ctx context.Context, uid uint, require bool) (*kolide.User, error) {
	user, err := svc.ds.UserByID(uid)
	if err != nil {
		return nil, errors.Wrap(err, "loading user by ID")
	}
	if user.SSOEnabled {
		return nil, errors.New("password reset for single sign on user not allowed")
	}
	// Require reset on next login
	user.AdminForcedPasswordReset = require
	if err := svc.saveUser(user); err != nil {
		return nil, errors.Wrap(err, "saving user")
	}

	if require {
		// Clear all of the existing sessions
		if err := svc.DeleteSessionsForUser(ctx, user.ID); err != nil {
			return nil, errors.Wrap(err, "deleting user sessions")
		}
	}

	return user, nil
}

func (svc service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := svc.ds.UserByEmail(email)
	if err != nil {
		return err
	}
	if user.SSOEnabled {
		return errors.New("password reset for single sign on user not allowed")
	}

	random, err := kolide.RandomText(svc.config.App.TokenKeySize)
	if err != nil {
		return err
	}
	token := base64.URLEncoding.EncodeToString([]byte(random))

	request := &kolide.PasswordResetRequest{
		ExpiresAt: time.Now().Add(time.Hour * 24),
		UserID:    user.ID,
		Token:     token,
	}
	request, err = svc.ds.NewPasswordResetRequest(request)
	if err != nil {
		return err
	}

	config, err := svc.AppConfig(ctx)
	if err != nil {
		return err
	}

	resetEmail := kolide.Email{
		Subject: "Reset Your Kolide Password",
		To:      []string{user.Email},
		Config:  config,
		Mailer: &kolide.PasswordResetMailer{
			KolideServerURL: template.URL(config.KolideServerURL),
			Token:           token,
		},
	}

	return svc.mailService.SendEmail(resetEmail)
}

// saves user in datastore.
// doesn't need to be exposed to the transport
// the service should expose actions for modifying a user instead
func (svc service) saveUser(user *kolide.User) error {
	return svc.ds.SaveUser(user)
}

// generateRandomText return a string generated by filling in keySize bytes with
// random data and then base64 encoding those bytes
func generateRandomText(keySize int) (string, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeEnableUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req enableUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeAdminUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req adminUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getUserRequest{ID: id}, nil
}

func decodeListUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listUsersRequest{ListOptions: opt}, nil
}

func decodeModifyUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeChangePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeRequirePasswordResetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, errors.Wrap(err, "getting ID from request")
	}

	var req requirePasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decoding JSON")
	}
	req.ID = id

	return req, nil
}

func decodePerformRequiredPasswordResetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req performRequiredPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decoding JSON")
	}
	return req, nil
}

func decodeForgotPasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req forgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeResetPasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type createUserRequest struct {
	payload kolide.UserPayload
}

type createUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r createUserResponse) error() error { return r.Err }

func makeCreateUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		user, err := svc.NewUser(ctx, req.payload)
		if err != nil {
			return createUserResponse{Err: err}, nil
		}
		return createUserResponse{User: user}, nil
	}
}

type getUserRequest struct {
	ID uint `json:"id"`
}

type getUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

func makeGetUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
		user, err := svc.User(ctx, req.ID)
		if err != nil {
			return getUserResponse{Err: err}, nil
		}
		return getUserResponse{User: user}, nil
	}
}

type adminUserRequest struct {
	ID    uint `json:"id"`
	Admin bool `json:"admin"`
}

type adminUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r adminUserResponse) error() error { return r.Err }

func makeAdminUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(adminUserRequest)
		user, err := svc.ChangeUserAdmin(ctx, req.ID, req.Admin)
		if err != nil {
			return adminUserResponse{Err: err}, nil
		}
		return adminUserResponse{User: user}, nil
	}
}

type enableUserRequest struct {
	ID      uint `json:"id"`
	Enabled bool `json:"enabled"`
}

type enableUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r enableUserResponse) error() error { return r.Err }

func makeEnableUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(enableUserRequest)
		user, err := svc.ChangeUserEnabled(ctx, req.ID, req.Enabled)
		if err != nil {
			return enableUserResponse{Err: err}, nil
		}
		return enableUserResponse{User: user}, nil
	}
}

func makeGetSessionUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, err := svc.AuthenticatedUser(ctx)
		if err != nil {
			return getUserResponse{Err: err}, nil
		}
		return getUserResponse{User: user}, nil
	}
}

type listUsersRequest struct {
	ListOptions kolide.ListOptions
}

type listUsersResponse struct {
	Users []kolide.User `json:"users"`
	Err   error         `json:"error,omitempty"`
}

func (r listUsersResponse) error() error { return r.Err }

func makeListUsersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listUsersRequest)
		users, err := svc.ListUsers(ctx, req.ListOptions)
		if err != nil {
			return listUsersResponse{Err: err}, nil
		}

		resp := listUsersResponse{Users: []kolide.User{}}
		for _, user := range users {
			resp.Users = append(resp.Users, *user)
		}
		return resp, nil
	}
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type changePasswordResponse struct {
	Err error `json:"error,omitempty"`
}

func (r changePasswordResponse) error() error { return r.Err }

func makeChangePasswordEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(changePasswordRequest)
		err := svc.ChangePassword(ctx, req.OldPassword, req.NewPassword)
		return changePasswordResponse{Err: err}, nil
	}
}

type resetPasswordRequest struct {
	PasswordResetToken string `json:"password_reset_token"`
	NewPassword        string `json:"new_password"`
}

type resetPasswordResponse struct {
	Err error `json:"error,omitempty"`
}

func (r resetPasswordResponse) error() error { return r.Err }

func makeResetPasswordEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(resetPasswordRequest)
		err := svc.ResetPassword(ctx, req.PasswordResetToken, req.NewPassword)
		return resetPasswordResponse{Err: err}, nil
	}
}

type modifyUserRequest struct {
	ID      uint
	payload kolide.UserPayload
}

type modifyUserResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r modifyUserResponse) error() error { return r.Err }

func makeModifyUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyUserRequest)
		user, err := svc.ModifyUser(ctx, req.ID, req.payload)
		if err != nil {
			return modifyUserResponse{Err: err}, nil
		}

		return modifyUserResponse{User: user}, nil
	}
}

type performRequiredPasswordResetRequest struct {
	Password string `json:"new_password"`
	ID       uint   `json:"id"`
}

type performRequiredPasswordResetResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r performRequiredPasswordResetResponse) error() error { return r.Err }

func makePerformRequiredPasswordResetEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(performRequiredPasswordResetRequest)
		user, err := svc.PerformRequiredPasswordReset(ctx, req.Password)
		if err != nil {
			return performRequiredPasswordResetResponse{Err: err}, nil
		}
		return performRequiredPasswordResetResponse{User: user}, nil
	}
}

type requirePasswordResetRequest struct {
	Require bool `json:"require"`
	ID      uint `json:"id"`
}

type requirePasswordResetResponse struct {
	User *kolide.User `json:"user,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r requirePasswordResetResponse) error() error { return r.Err }

func makeRequirePasswordResetEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(requirePasswordResetRequest)
		user, err := svc.RequirePasswordReset(ctx, req.ID, req.Require)
		if err != nil {
			return requirePasswordResetResponse{Err: err}, nil
		}
		return requirePasswordResetResponse{User: user}, nil
	}
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type forgotPasswordResponse struct {
	Err error `json:"error,omitempty"`
}

func (r forgotPasswordResponse) error() error { return r.Err }
func (r forgotPasswordResponse) status() int  { return http.StatusAccepted }

func makeForgotPasswordEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(forgotPasswordRequest)
		err := svc.RequestPasswordReset(ctx, req.Email)
		if err != nil {
			return forgotPasswordResponse{Err: err}, nil
		}
		return forgotPasswordResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) ChangeUserAdmin(ctx context.Context, id uint, isAdmin bool) (*kolide.User, error) {
	var (
		user *kolide.User
		err  error
	)

	defer func(begin time.Time) {
		lvs := []string{"method", "ChangeUserAdmin", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.Service.ChangeUserAdmin(ctx, id, isAdmin)
	return user, err
}

func (mw metricsMiddleware) ChangeUserEnabled(ctx context.Context, id uint, isEnabled bool) (*kolide.User, error) {
	var (
		user *kolide.User
		err  error
	)

	defer func(begin time.Time) {
		lvs := []string{"method", "ChangeUserEnabled", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.Service.ChangeUserEnabled(ctx, id, isEnabled)
	return user, err
}

func (mw metricsMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	var (
		user *kolide.User
		err  error
	)

	defer func(begin time.Time) {
		lvs := []string{"method", "NewUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.Service.NewUser(ctx, p)
	return user, err
}

func (mw metricsMiddleware) ModifyUser(ctx context.Context, userID uint, p kolide.UserPayload) (*kolide.User, error) {
	var (
		user *kolide.User
		err  error
	)

	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.Service.ModifyUser(ctx, userID, p)
	return user, err
}

func (mw metricsMiddleware) User(ctx context.Context, id uint) (*kolide.User, error) {
	var (
		user *kolide.User
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "User", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.Service.User(ctx, id)
	return user, err
}

func (mw metricsMiddleware) ListUsers(ctx context.Context, opt kolide.ListOptions) ([]*kolide.User, error) {

	var (
		users []*kolide.User
		err   error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "Users", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	users, err = mw.Service.ListUsers(ctx, opt)
	return users, err
}

func (mw metricsMiddleware) AuthenticatedUser(ctx context.Context) (*kolide.User, error) {
	var (
		user *kolide.User
		err  error
	)

	defer func(begin time.Time) {
		lvs := []string{"method", "AuthenticatedUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.Service.AuthenticatedUser(ctx)
	return user, err
}

func (mw metricsMiddleware) ChangePassword(ctx context.Context, oldPass, newPass string) error {
	var err error

	defer func(begin time.Time) {
		lvs := []string{"method", "ChangePassword", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Service.ChangePassword(ctx, oldPass, newPass)
	return err
}

func (mw metricsMiddleware) ResetPassword(ctx context.Context, token, password string) error {
	var err error

	defer func(begin time.Time) {
		lvs := []string{"method", "ResetPassword", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Service.ResetPassword(ctx, token, password)
	return err
}

func (mw metricsMiddleware) RequestPasswordReset(ctx context.Context, email string) error {
	var err error

	defer func(begin time.Time) {
		lvs := []string{"method", "RequestPasswordReset", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Service.RequestPasswordReset(ctx, email)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) ChangeUserAdmin(ctx context.Context, id uint, isAdmin bool) (*kolide.User, error) {
	var (
		loggedInUser = "unauthenticated"
		userName     = "none"
		err          error
		user         *kolide.User
	)

	vc, ok := viewer.FromContext(ctx)
	if ok {
		loggedInUser = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ChangeUserAdmin",
			"user", userName,
			"changed_by", loggedInUser,
			"admin", isAdmin,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.ChangeUserAdmin(ctx, id, isAdmin)
	if user != nil {
		userName = user.Username
	}
	return user, err
}

func (mw loggingMiddleware) ChangeUserEnabled(ctx context.Context, id uint, isEnabled bool) (*kolide.User, error) {
	var (
		loggedInUser = "unauthenticated"
		userName     = "none"
		err          error
		user         *kolide.User
	)

	vc, ok := viewer.FromContext(ctx)
	if ok {
		loggedInUser = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ChangeUserEnabled",
			"user", userName,
			"changed_by", loggedInUser,
			"enabled", isEnabled,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.ChangeUserEnabled(ctx, id, isEnabled)
	if user != nil {
		userName = user.Username
	}
	return user, err
}

func (mw loggingMiddleware) NewAdminCreatedUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	var (
		user         *kolide.User
		err          error
		username     = "none"
		loggedInUser = "unauthenticated"
	)

	vc, ok := viewer.FromContext(ctx)
	if ok {
		loggedInUser = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "NewAdminCreatedUser",
			"user", username,
			"created_by", loggedInUser,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.NewAdminCreatedUser(ctx, p)
	if user != nil {
		username = user.Username
	}
	return user, err
}

func (mw loggingMiddleware) ListUsers(ctx context.Context, opt kolide.ListOptions) ([]*kolide.User, error) {
	var (
		users    []*kolide.User
		err      error
		username = "none"
	)

	vc, ok := viewer.FromContext(ctx)
	if ok {
		username = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListUsers",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	users, err = mw.Service.ListUsers(ctx, opt)
	return users, err
}

func (mw loggingMiddleware) RequirePasswordReset(ctx context.Context, uid uint, require bool) (*kolide.User, error) {
	var (
		user     *kolide.User
		err      error
		username = "none"
	)

	vc, ok := viewer.FromContext(ctx)
	if ok {
		username = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "RequirePasswordReset",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.RequirePasswordReset(ctx, uid, require)
	return user, err

}

func (mw loggingMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	var (
		user         *kolide.User
		err          error
		username     = "none"
		loggedInUser = "unauthenticated"
	)

	vc, ok := viewer.FromContext(ctx)
	if ok {
		loggedInUser = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "NewUser",
			"user", username,
			"created_by", loggedInUser,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.NewUser(ctx, p)

	if user != nil {
		username = user.Username
	}
	return user, err
}

func (mw loggingMiddleware) ModifyUser(ctx context.Context, userID uint, p kolide.UserPayload) (*kolide.User, error) {
	var (
		user     *kolide.User
		err      error
		username = "none"
	)

	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ModifyUser",
			"user", username,
			"modified_by", vc.Username(),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.ModifyUser(ctx, userID, p)

	if user != nil {
		username = user.Username
	}

	return user, err
}

func (mw loggingMiddleware) User(ctx context.Context, id uint) (*kolide.User, error) {
	var (
		user     *kolide.User
		err      error
		username = "none"
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "User",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.User(ctx, id)

	if user != nil {
		username = user.Username
	}
	return user, err
}

func (mw loggingMiddleware) AuthenticatedUser(ctx context.Context) (*kolide.User, error) {
	var (
		user     *kolide.User
		err      error
		username = "none"
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "User",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.AuthenticatedUser(ctx)

	if user != nil {
		username = user.Username
	}
	return user, err
}

func (mw loggingMiddleware) ChangePassword(ctx context.Context, oldPass, newPass string) error {
	var (
		requestedBy = "unauthenticated"
		err         error
	)
	vc, ok := viewer.FromContext(ctx)
	if ok {
		requestedBy = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ChangePassword",
			"err", err,
			"requested_by", requestedBy,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.ChangePassword(ctx, oldPass, newPass)
	return err
}

func (mw loggingMiddleware) ResetPassword(ctx context.Context, token, password string) error {
	var err error

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ResetPassword",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.ResetPassword(ctx, token, password)
	return err
}

func (mw loggingMiddleware) RequestPasswordReset(ctx context.Context, email string) error {
	var (
		requestedBy = "unauthenticated"
		err         error
	)
	vc, ok := viewer.FromContext(ctx)
	if ok {
		requestedBy = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "RequestPasswordReset",
			"email", email,
			"err", err,
			"requested_by", requestedBy,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.RequestPasswordReset(ctx, email)
	return err
}

func (mw loggingMiddleware) PerformRequiredPasswordReset(ctx context.Context, password string) (*kolide.User, error) {
	var (
		resetBy = "unauthenticated"
		err     error
	)
	vc, ok := viewer.FromContext(ctx)
	if ok {
		resetBy = vc.Username()
	}
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "PerformRequiredPasswordReset",
			"err", err,
			"reset_by", resetBy,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err := mw.Service.PerformRequiredPasswordReset(ctx, password)
	return user, err
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func (mw validationMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	invalid := &invalidArgumentError{}
	if p.Username == nil {
		invalid.Append("username", "missing required argument")
	} else {
		if *p.Username == "" {
			invalid.Append("username", "cannot be empty")
		}

		if strings.Contains(*p.Username, "@") {
			invalid.Append("username", "'@' character not allowed in usernames")
		}
	}

	// we don't need a password for single sign on
	if p.SSOInvite == nil || !*p.SSOInvite {
		if p.Password == nil {
			invalid.Append("password", "missing required argument")
		} else {
			if *p.Password == "" {
				invalid.Append("password", "cannot be empty")
			}
			if err := validatePasswordRequirements(*p.Password); err != nil {
				invalid.Append("password", err.Error())
			}
		}
	}

	if p.Email == nil {
		invalid.Append("email", "missing required argument")
	} else {
		if *p.Email == "" {
			invalid.Append("email", "cannot be empty")
		}
	}

	if p.InviteToken == nil {
		invalid.Append("invite_token", "missing required argument")
	} else {
		if *p.InviteToken == "" {
			invalid.Append("invite_token", "cannot be empty")
		}
	}

	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.NewUser(ctx, p)
}

func (mw validationMiddleware) ModifyUser(ctx context.Context, userID uint, p kolide.UserPayload) (*kolide.User, error) {
	invalid := &invalidArgumentError{}
	if p.Username != nil {
		if *p.Username == "" {
			invalid.Append("username", "cannot be empty")
		}

		if strings.Contains(*p.Username, "@") {
			invalid.Append("username", "'@' character not allowed in usernames")
		}
	}

	if p.Name != nil {
		if *p.Name == "" {
			invalid.Append("name", "cannot be empty")
		}
	}

	if p.Email != nil {
		if *p.Email == "" {
			invalid.Append("email", "cannot be empty")
		}
		// if the user is not an admin, or if an admin is changing their own email
		// address a password is required,
		if passwordRequiredForEmailChange(ctx, userID, invalid) {
			if p.Password == nil {
				invalid.Append("password", "cannot be empty if email is changed")
			}
		}
	}

	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyUser(ctx, userID, p)
}

func passwordRequiredForEmailChange(ctx context.Context, uid uint, invalid *invalidArgumentError) bool {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		invalid.Append("viewer", "not present")
		return false
	}
	// if a user is changing own email need a password no matter what
	if vc.UserID() == uid {
		return true
	}
	// if an admin is changing another users email no password needed
	if vc.IsAdmin() {
		return false
	}
	// should never get here because a non admin can't change the email of another
	// user
	invalid.Append("auth", "this user can't change another user's email")
	return false
}

func (mw validationMiddleware) ChangePassword(ctx context.Context, oldPass, newPass string) error {
	invalid := &invalidArgumentError{}
	if oldPass == "" {
		invalid.Append("old_password", "cannot be empty")
	}
	if newPass == "" {
		invalid.Append("new_password", "cannot be empty")
	}

	if err := validatePasswordRequirements(newPass); err != nil {
		invalid.Append("new_password", err.Error())
	}

	if invalid.HasErrors() {
		return invalid
	}
	return mw.Service.ChangePassword(ctx, oldPass, newPass)
}

func (mw validationMiddleware) ResetPassword(ctx context.Context, token, password string) error {
	invalid := &invalidArgumentError{}
	if token == "" {
		invalid.Append("token", "cannot be empty field")
	}
	if password == "" {
		invalid.Append("new_password", "cannot be empty field")
	}
	if err := validatePasswordRequirements(password); err != nil {
		invalid.Append("new_password", err.Error())
	}
	if invalid.HasErrors() {
		return invalid
	}
	return mw.Service.ResetPassword(ctx, token, password)
}

// Requirements for user password:
// at least 7 character length
// at least 1 symbol
// at least 1 number
func validatePasswordRequirements(password string) error {
	var (
		number bool
		symbol bool
	)

	for _, s := range password {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			symbol = true
		}
	}

	if len(password) >= 7 &&
		number &&
		symbol {
		return nil
	}

	return errors.New("password does not meet validation requirements")
}
