package app

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/errors"
	"github.com/kolide/kolide-ose/kolide"
	"github.com/kolide/kolide-ose/sessions"
)

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
// User management web endpoints
////////////////////////////////////////////////////////////////////////////////

// swagger:parameters GetUser
type GetUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// swagger:response GetUserResponseBody
type GetUserResponseBody struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Admin              bool   `json:"admin"`
	Enabled            bool   `json:"enabled"`
	NeedsPasswordReset bool   `json:"needs_password_reset"`
}

// swagger:route POST /api/v1/kolide/user GetUser
//
// Get information on a user
//
// Using this API will allow the requester to inspect and get info on
// other users in the application.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func GetUser(c *gin.Context) {
	var body GetUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.IsLoggedIn() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformReadActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters CreateUser
type CreateUserRequestBody struct {
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Admin              bool   `json:"admin"`
	NeedsPasswordReset bool   `json:"needs_password_reset"`
}

// swagger:route PUT /api/v1/kolide/user CreateUser
//
// Create a new user
//
// Using this API will allow the requester to create a new user with the ability
// to control various user settings
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func CreateUser(c *gin.Context) {
	var body CreateUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	// TODO move this to middleware
	vc := VC(c)
	if !vc.IsAdmin() {
		UnauthorizedError(c)
		return
	}

	// temporary, pass args explicitly as well
	db := GetDB(c)

	u, err := kolide.NewUser(body.Username, body.Password, body.Email, body.Admin, body.NeedsPasswordReset)
	if err != nil {
		logrus.Errorf("Error creating new user: %s", err.Error())
		errors.ReturnError(c, err)
		return
	}

	// save user in db
	user, err := db.NewUser(u)
	if err != nil {
		logrus.Errorf("error creating new user: %s", err)
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters ModifyUser
type ModifyUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// swagger:route PATCH /api/v1/kolide/user ModifyUser
//
// Update a user's basic information and settings
//
// Using this API will allow the requester to update a user's basic settings.
// Note that updating administrative settings are not exposed via this endpoint
// as this is primarily intended to be used by users to update their own
// settings.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func ModifyUser(c *gin.Context) {
	var body ModifyUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	if body.Name != "" {
		user.Name = body.Name
	}
	if body.Email != "" {
		user.Email = body.Email
	}
	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}
	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters ChangeUserPassword
type ChangePasswordRequestBody struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	CurrentPassword    string `json:"current_password"`
	PasswordResetToken string `json:"password_reset_token"`
	NewPassword        string `json:"new_password" validate:"required"`
	NewPasswordConfim  string `json:"new_password_confirm" validate:"required"`
}

// swagger:route PATCH /api/v1/kolide/user/password ChangeUserPassword
//
// Change a user's password
//
// Using this API will allow the requester to change their password. Users
// should include their own user id as the "id" paramater and/or their own
// username as the "username" parameter. Admins can change the passords for
// other users by defining their ID or username.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func ChangeUserPassword(c *gin.Context) {
	var body ChangePasswordRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	if body.NewPassword != body.NewPasswordConfim {
		c.JSON(406, map[string]interface{}{"error": "Passwords do not match"})
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	deleteResetRequest := func(reset *kolide.PasswordResetRequest) {
		if err != nil { //TODO: this shadows error returns
			err = db.DeletePasswordResetRequest(reset)
			if err != nil {
				errors.ReturnError(c, errors.DatabaseError(err))
				return
			}
		}
	}

	if body.PasswordResetToken != "" {
		reset, err := db.FindPassswordResetByToken(body.PasswordResetToken)
		if err != nil {
			UnauthorizedError(c)
			return
		}

		if time.Now().After(reset.ExpiresAt) {
			deleteResetRequest(reset)
			UnauthorizedError(c)
			return
		}
		defer deleteResetRequest(reset)
	} else if !vc.IsAdmin() {
		if body.CurrentPassword != "" {
			if user.ValidatePassword(body.CurrentPassword) != nil {
				UnauthorizedError(c)
				return
			}
		} else {
			UnauthorizedError(c)
			return
		}
	}

	err = user.SetPassword(body.NewPassword)
	if err != nil {
		logrus.Errorf("Error setting user password: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err)) // probably not this
		return
	}

	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters SetUserAdminState
type SetUserAdminStateRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

// swagger:route PATCH /api/v1/kolide/user/admin SetUserAdminState
//
// Modify a user's admin settings
//
// This endpoint allows an existing admin to promote a non-admin to admin or
// demote a current admin to non-admin.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func SetUserAdminState(c *gin.Context) {
	var body SetUserAdminStateRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.IsAdmin() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	user.Admin = body.Admin
	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}
	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters SetUserEnabledState
type SetUserEnabledStateRequestBody struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	Enabled         bool   `json:"enabled"`
	CurrentPassword string `json:"current_password"`
}

// swagger:route PATCH /api/v1/kolide/user/enabled SetUserEnabledState
//
// Enable or disable a user.
//
// This endpoint allows an existing admin to enable a disabled user or
// disable an enabled user. If a user calls this endpoint, to disable,
// their own account, they must also submit their current password, to
// verify their request.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func SetUserEnabledState(c *gin.Context) {
	var body SetUserEnabledStateRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	if !vc.IsAdmin() {
		if user.ValidatePassword(body.CurrentPassword) != nil {
			UnauthorizedError(c)
			return
		}
	}

	user.Enabled = body.Enabled
	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}
	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

///////////////////////////////////////////////////////////////////////////////
// Session management HTTP endpoints
////////////////////////////////////////////////////////////////////////////////

// Setting the session backend via a middleware
func SessionBackendMiddleware(backend sessions.SessionBackend) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("SessionBackend", backend)
		c.Next()
	}
}

// Get the database connection from the context, or panic
func GetSessionBackend(c *gin.Context) sessions.SessionBackend {
	return c.MustGet("SessionBackend").(sessions.SessionBackend)
}

////////////////////////////////////////////////////////////////////////////////
// Session management HTTP endpoints
////////////////////////////////////////////////////////////////////////////////

// swagger:parameters DeleteSession
type DeleteSessionRequestBody struct {
	SessionID uint `json:"session_id" validate:"required"`
}

// swagger:route DELETE /api/v1/kolide/session DeleteSession
//
// Delete a specific session, as specified by the session's ID.
//
// This API allows for the requester to delete a specific session. Note that the
// API expects the session ID as the parameter, not the session key.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: nil
func DeleteSession(c *gin.Context) {
	var body DeleteSessionRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	sb := GetSessionBackend(c)

	session, err := sb.FindID(body.SessionID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(session.UserID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	err = sb.Destroy(session)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// swagger:parameters DeleteSessionsForUser
type DeleteSessionsForUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// swagger:route DELETE /api/v1/kolide/user/sessions DeleteSessionsForUser
//
// Delete all of a user's active sessions
//
// This API allows an admin to delete all active sessions that are known to
// belong to a specific user. This effectively logs out the user on all
// devices.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: nil
func DeleteSessionsForUser(c *gin.Context) {
	var body DeleteSessionsForUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	sb := GetSessionBackend(c)
	err = sb.DestroyAllForUser(user.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, nil)

}

// swagger:parameters GetInfoAboutSession
type GetInfoAboutSessionRequestBody struct {
	SessionKey string `json:"session_key" validate:"required"`
}

// swagger:response SessionInfoResponseBody
type SessionInfoResponseBody struct {
	SessionID  uint      `json:"session_id"`
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	AccessedAt time.Time `json:"created_at"`
}

// swagger:route POST /api/v1/kolide/session GetInfoAboutSession
//
// Get information on a session, given a session key.
//
// Using this API will allow the requester to inspect and get info on
// an active session, given the session key.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: SessionInfoResponseBody
func GetInfoAboutSession(c *gin.Context) {
	var body GetInfoAboutSessionRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	sb := GetSessionBackend(c)
	session, err := sb.FindKey(body.SessionKey)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(session.UserID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.IsAdmin() && !vc.IsUserID(user.ID) {
		UnauthorizedError(c)
		return
	}

	c.JSON(http.StatusOK, &SessionInfoResponseBody{
		SessionID:  session.ID,
		UserID:     session.UserID,
		CreatedAt:  session.CreatedAt,
		AccessedAt: session.AccessedAt,
	})
}

// swagger:parameters GetInfoAboutSessionsForUser
type GetInfoAboutSessionsForUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// swagger:response GetInfoAboutSessionsForUserResponseBody
type GetInfoAboutSessionsForUserResponseBody struct {
	Sessions []SessionInfoResponseBody `json:"sessions"`
}

// swagger:route POST /api/v1/kolide/user/sessions GetInfoAboutSessionsForUser
//
// Get information on a user's sessions
//
// Using this API will allow an admin to inspect and get info on all of a user's
// active session.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetInfoAboutSessionsForUserResponseBody
func GetInfoAboutSessionsForUser(c *gin.Context) {
	var body GetInfoAboutSessionsForUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	sb := GetSessionBackend(c)
	sessions, err := sb.FindAllForUser(user.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	var response []SessionInfoResponseBody
	for _, session := range sessions {
		response = append(response, SessionInfoResponseBody{
			SessionID:  session.ID,
			UserID:     session.UserID,
			CreatedAt:  session.CreatedAt,
			AccessedAt: session.AccessedAt,
		})
	}

	c.JSON(http.StatusOK, &GetInfoAboutSessionsForUserResponseBody{
		Sessions: response,
	})
}

////////////////////////////////////////////////////////////////////////////////
// Password Reset HTTP endpoints
////////////////////////////////////////////////////////////////////////////////

// swagger:parameters ResetUserPassword
type ResetPasswordRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// swagger:response ResetPasswordResponseBody
type ResetPasswordResponseBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// swagger:route POST /api/v1/kolide/user/password/reset ResetUserPassword
//
// Reset a user's password
//
// Using this API will allow the requester to reset their password. Users
// should include their own user id as the "id" paramater and/or their own
// username as the "username" parameter. Admins can change the passwords for
// other users by defining their ID or username. Logged out users can reset
// their own password by including their email in addition to either their
// user id or their username.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: ResetPasswordResponseBody
func ResetUserPassword(c *gin.Context) {
	var body ResetPasswordRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		logrus.Errorf("Error parsing ResetPassword post body: %s", err.Error())
		return
	}

	db := GetDB(c)
	vc := VC(c)

	var userEmail string
	if !vc.IsLoggedIn() {
		if body.Email == "" {
			UnauthorizedError(c)
			return
		}
		userEmail = body.Email
	}

	// TODO: Should this find user by email if user isn't logged in?
	// gorm makes this confusing...
	_ = userEmail
	user, err := db.UserByID(body.ID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			NotFoundRequestError(c)
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	if vc.IsAdmin() || vc.IsUserID(user.ID) || !vc.IsLoggedIn() {
		// logged-in admin user resetting a password or logged-in user
		// resetting their own password or logged-out user presumably resetting
		// their own password

		if vc.CanPerformWriteActionOnUser(user) {
			// if the user is logged out, don't perform the user state
			// modifications
			user.NeedsPasswordReset = true

			err = db.SaveUser(user)
			if err != nil {
				errors.ReturnError(c, errors.DatabaseError(err))
				return
			}
		}

		request, err := kolide.NewPasswordResetRequest(GetDB(c), user.ID, time.Now().Add(time.Hour*24))
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Database error"))
			return
		}

		html, text, err := GetEmailBody(PasswordResetEmail, PasswordResetRequestEmailParameters{
			Name:  user.Name,
			Token: request.Token,
		})
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Email error"))
			return
		}

		subject, err := GetEmailSubject(PasswordResetEmail)
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Email error"))
			return
		}

		err = SendEmail(GetSMTPConnectionPool(c), user.Email, subject, html, text)
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Email error"))
			return
		}
	} else {
		// Logged-in user trying to reset another user's password
		UnauthorizedError(c)
		return
	}

	c.JSON(http.StatusOK, ResetPasswordResponseBody{
		ID:       user.ID,
		Username: user.Username,
	})
}

// swagger:parameters VerifyPasswordResetRequest
type VerifyPasswordResetRequestRequestBody struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Token    string `json:"token"`
}

// swagger:parameters VerifyPasswordResetRequestResponseBody
type VerifyPasswordResetRequestResponseBody struct {
	Valid bool `json:"valid"`
	ID    uint `json:"id"`
}

// swagger:route POST /api/v1/kolide/user/password/verify VerifyPasswordResetRequest
//
// Verify an email campaign before it is used.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: VerifyPasswordResetRequestResponseBody
func VerifyPasswordResetRequest(c *gin.Context) {
	var body VerifyPasswordResetRequestRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		logrus.Errorf("Error parsing request body: %s", err.Error())
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.UserID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, &VerifyPasswordResetRequestResponseBody{
				Valid: false,
			})
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	reset, err := db.FindPassswordResetByTokenAndUserID(body.Token, user.ID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, VerifyPasswordResetRequestResponseBody{
				Valid: false,
			})
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	if time.Now().After(reset.ExpiresAt) {
		c.JSON(http.StatusNotFound, VerifyPasswordResetRequestResponseBody{
			Valid: false,
		})
		return
	}

	c.JSON(http.StatusOK, VerifyPasswordResetRequestResponseBody{
		Valid: true,
		ID:    reset.ID,
	})
}

// swagger:parameters DeletePasswordResetRequest
type DeletePasswordResetRequestRequestBody struct {
	ID uint `json:"id" validate:"required"`
}

// swagger:route DELETE /api/v1/kolide/user/password/reset DeletePasswordResetRequest
//
// Delete an email campaign after it has been used.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: nil
func DeletePasswordResetRequest(c *gin.Context) {
	var body DeletePasswordResetRequestRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		logrus.Errorf("Error parsing request body: %s", err.Error())
		return
	}

	db := GetDB(c)
	campaign, err := db.FindPassswordResetByID(body.ID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			NotFoundRequestError(c)
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	user, err := db.UserByID(campaign.UserID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	vc := VC(c)
	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	err = db.DeletePasswordResetRequest(campaign)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}
