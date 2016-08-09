package app

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// ViewerContext is a struct which represents the ability for an execution
// context to participate in certain actions. Most often, a ViewerContext is
// associated with an application user, but a ViewerContext can represent a
// variety of other execution contexts as well (script, test, etc). The main
// purpose of a ViewerContext is to assist in the authorization of sensitive
// actions.
type ViewerContext struct {
	user *User
}

// IsAdmin indicates whether or not the current user can perform administrative
// actions.
func (vc *ViewerContext) IsAdmin() bool {
	if vc.user != nil {
		return vc.user.Admin && vc.user.Enabled
	}
	return false
}

// UserID is a helper that enables quick access to the user ID of the current
// user.
func (vc *ViewerContext) UserID() (uint, error) {
	if vc.user != nil {
		return vc.user.ID, nil
	}
	return 0, errors.New("Unauthorized", "No user set")
}

// CanPerformActions returns a bool indicating the current user's ability to
// perform the most basic actions on the site
func (vc *ViewerContext) CanPerformActions() bool {
	if vc.IsLoggedIn() && !vc.user.NeedsPasswordReset {
		return true
	}

	return false
}

// IsLoggedIn determines whether or not the current VC is attached to a user
// account
func (vc *ViewerContext) IsLoggedIn() bool {
	if vc.user != nil && vc.user.Enabled {
		return true
	}

	return false
}

// IsUserID returns true if the given user id the same as the user which is
// represented by this ViewerContext
func (vc *ViewerContext) IsUserID(id uint) bool {
	userID, err := vc.UserID()
	if err != nil {
		return false
	}
	if userID == id {
		return true
	}
	return false
}

// CanPerformWriteActionsOnUser returns a bool indicating the current user's
// ability to perform write actions on the given user
func (vc *ViewerContext) CanPerformWriteActionOnUser(u *User) bool {
	return vc.CanPerformActions() && (vc.IsUserID(u.ID) || vc.IsAdmin())
}

// CanPerformReadActionsOnUser returns a bool indicating the current user's
// ability to perform read actions on the given user
func (vc *ViewerContext) CanPerformReadActionOnUser(u *User) bool {
	return vc.CanPerformActions() || (vc.IsLoggedIn() && vc.IsUserID(u.ID))
}

// GenerateVC generates a ViewerContext given a user struct
func GenerateVC(user *User) *ViewerContext {
	return &ViewerContext{
		user: user,
	}
}

// EmptyVC is a utility which generates an empty ViewerContext. This is often
// used to represent users which are not logged in.
func EmptyVC() *ViewerContext {
	return &ViewerContext{
		user: nil,
	}
}

// VC accepts a web request context and a database handler and attempts
// to parse a user's jwt token out of the active session, validate the token,
// and generate an appropriate ViewerContext given the data in the session.
func VC(c *gin.Context) *ViewerContext {
	sm := NewSessionManager(c)
	session, err := sm.Session()
	if err != nil {
		return EmptyVC()
	}
	return VCForID(GetDB(c), session.UserID)
}

func VCForID(db *gorm.DB, id uint) *ViewerContext {
	// Generating a VC requires a user struct. Attempt to populate one using
	// the user id of the current session holder
	user := &User{ID: id}
	err := db.Where(user).First(user).Error
	if err != nil {
		return EmptyVC()
	}

	return GenerateVC(user)
}

////////////////////////////////////////////////////////////////////////////////
// Login and password utilities
////////////////////////////////////////////////////////////////////////////////

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

func HashPassword(salt, password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%s%s", password, salt)),
		viper.GetInt("auth.bcrypt_cost"),
	)
}

func SaltAndHashPassword(password string) (string, []byte, error) {
	salt, err := generateRandomText(viper.GetInt("auth.salt_key_size"))
	if err != nil {
		return "", []byte{}, err
	}
	hashed, err := HashPassword(salt, password)
	return salt, hashed, err
}

////////////////////////////////////////////////////////////////////////////////
// Authentication and authorization web endpoints
////////////////////////////////////////////////////////////////////////////////

// swagger:parameters Login
type LoginRequestBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// swagger:route POST /api/v1/kolide/login Login
//
// Login to the application
//
// This allows you to submit a set of credentials to the server and the
// server will attempt to validate the credentials against the backing
// database. If the credentials are valid, the response will issue the
// browser a new session cookie.
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
//       authenticated: no
//
//     Responses:
//       200: GetUserResponseBody
func Login(c *gin.Context) {
	var body LoginRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	db := GetDB(c)

	user := &User{Username: body.Username}
	err = db.Where(user).First(user).Error
	if err != nil {
		logrus.Debugf("User not found: %s", body.Username)
		UnauthorizedError(c)
		return
	}

	err = user.ValidatePassword(body.Password)
	if err != nil {
		logrus.Debugf("Invalid password for user: %s", body.Username)
		UnauthorizedError(c)
		return
	}

	sm := NewSessionManager(c)
	sm.MakeSessionForUserID(user.ID)
	err = sm.Save()
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(200, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:route GET /api/v1/kolide/logout Logout
//
// Logout of the application
//
// This endpoint will delete the current session from the backend database
// and log the user out of the application
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
func Logout(c *gin.Context) {
	sm := NewSessionManager(c)

	err := sm.Destroy()
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	err = sm.Save()
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(200, nil)
}
