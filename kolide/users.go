package kolide

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// UserStore contains methods for managing users in a datastore
type UserStore interface {
	NewUser(user *User) (*User, error)
	User(username string) (*User, error)
	UserByID(id uint) (*User, error)
	SaveUser(user *User) error
}

type UserService interface {
	NewUser(ctx context.Context, p UserPayload) (*User, error)
	User(ctx context.Context, id uint) (*User, error)
	ChangePassword(ctx context.Context, userID uint, old, new string) error
	UpdateAdminRole(ctx context.Context, userID uint, isAdmin bool) error
	UpdateUserStatus(ctx context.Context, userID uint, password string, enabled bool) error
}

// User is the model struct which represents a kolide user
type User struct {
	ID                       uint `gorm:"primary_key"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
	Username                 string `gorm:"not null;unique_index:idx_user_unique_username"`
	Password                 []byte `gorm:"not null"`
	Salt                     string `gorm:"not null"`
	Name                     string
	Email                    string `gorm:"not null;unique_index:idx_user_unique_email"`
	Admin                    bool   `gorm:"not null"`
	Enabled                  bool   `gorm:"not null"`
	AdminForcedPasswordReset bool
}

// UserPayload is used to modify an existing user
type UserPayload struct {
	Username                 *string `json:"username"`
	Name                     *string `json:"name"`
	Email                    *string `json:"email"`
	Admin                    *bool   `json:"admin"`
	Enabled                  *bool   `json:"enabled"`
	AdminForcedPasswordReset *bool   `json:"force_password_reset"`
	Password                 *string `json:"password"`
	// modify params
	CurrentPassword *string `json:"current_password"`
	NewPassword     *string `json:"new_password"`
}

// NewUser is a wrapper around the creation of a new user.
// NewUser exists largely to allow the API to simply accept a string password
// while using the applications password hashing mechanisms to salt and hash the
// password.
func NewUser(username, password, email string, admin, needsPasswordReset bool) (*User, error) {
	salt, hash, err := saltAndHashPassword(password)
	if err != nil {
		return nil, err
	}
	user := User{
		Username:                 username,
		Password:                 hash,
		Salt:                     salt,
		Email:                    email,
		Admin:                    admin,
		Enabled:                  true,
		AdminForcedPasswordReset: needsPasswordReset,
	}
	return &user, nil
}

// ValidatePassword accepts a potential password for a given user and attempts
// to validate it against the hash stored in the database after joining the
// supplied password with the stored password salt
func (u *User) ValidatePassword(password string) error {
	saltAndPass := []byte(fmt.Sprintf("%s%s", password, u.Salt))
	return bcrypt.CompareHashAndPassword(u.Password, saltAndPass)
}

// SetPassword accepts a new password for a user object and updates the salt
// and hash for that user in the database. This function assumes that the
// appropriate authorization checks have already occurred by the caller.
func (u *User) SetPassword(password string) error {
	salt, hash, err := saltAndHashPassword(password)
	if err != nil {
		return err
	}
	u.Salt = salt
	u.Password = hash
	return nil
}

// TODO make viper config specific
func hashPassword(salt, password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%s%s", password, salt)),
		viper.GetInt("auth.bcrypt_cost"),
	)
}

// TODO make viper config specific
func saltAndHashPassword(password string) (string, []byte, error) {
	salt, err := generateRandomText(viper.GetInt("auth.salt_key_size"))
	if err != nil {
		return "", []byte{}, err
	}
	hashed, err := hashPassword(salt, password)
	return salt, hashed, err
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
