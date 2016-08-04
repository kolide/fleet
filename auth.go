package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	return 0, errors.New("No user set")
}

// CanPerformActions returns a bool indicating the current user's ability to
// perform the most basic actions on the site
func (vc *ViewerContext) CanPerformActions() bool {
	if vc.user == nil {
		return false
	}

	if !vc.user.Enabled {
		return false
	}

	return true
}

// IsUserID return true if the given user id the same as the user which is
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
	return vc.CanPerformActions()
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
func VC(c *gin.Context, db *gorm.DB) (*ViewerContext, error) {
	sm := NewSessionManager(c.Request, c.Writer, &GormSessionBackend{db: db}, db)
	vc := sm.VC()
	return vc, nil
}

////////////////////////////////////////////////////////////////////////////////
// JSON Web Tokens
////////////////////////////////////////////////////////////////////////////////

// Given a session key create a JWT to be delivered to the client
func GenerateJWT(sessionKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_key": sessionKey,
	})

	return token.SignedString([]byte(config.App.JWTKey))
}

// ParseJWT attempts to parse a JWT token in serialized string form into a
// JWT token in a deserialized jwt.Token struct.
func ParseJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		method, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok || method != jwt.SigningMethodHS256 {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(config.App.JWTKey), nil
	})
}

////////////////////////////////////////////////////////////////////////////////
// Login and password utilities
////////////////////////////////////////////////////////////////////////////////

func generateRandomText(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func HashPassword(salt, password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%s%s", salt, password)),
		config.App.BcryptCost,
	)
}

func SaltAndHashPassword(password string) (string, []byte, error) {
	salt := generateRandomText(config.App.SaltLength)
	hashed, err := HashPassword(salt, password)
	return salt, hashed, err
}

////////////////////////////////////////////////////////////////////////////////
// Authentication and authorization web endpoints
////////////////////////////////////////////////////////////////////////////////

type LoginRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var body LoginRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Errorf("Error parsing Login post body: %s", err.Error())
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

	sm := NewSessionManager(c.Request, c.Writer, &GormSessionBackend{db: db}, db)
	sm.MakeSessionForUser(user)
	err = sm.Save()
	if err != nil {
		DatabaseError(c)
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

func Logout(c *gin.Context) {
	db := GetDB(c)
	sm := NewSessionManager(c.Request, c.Writer, &GormSessionBackend{db: db}, db)

	err := sm.Destroy()
	if err != nil {
		DatabaseError(c)
		return
	}

	err = sm.Save()
	if err != nil {
		DatabaseError(c)
		return
	}

	c.JSON(200, nil)
}
