package kolide

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// UserStore contains methods for managing users in a datastore
type UserStore interface {
	NewUser(user *User) (*User, error)
	User(username string) (*User, error)
	UserByEmail(email string) (*User, error)
	UserByID(id uint) (*User, error)
	SaveUser(user *User) error
}

type UserService interface {
	NewUser(ctx context.Context, p UserPayload) (*User, error)
	User(ctx context.Context, id uint) (*User, error)
	ChangePassword(ctx context.Context, userID uint, old, new string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ModifyUser(ctx context.Context, userID uint, p UserPayload) (*User, error)
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
	GravatarURL              string
	Position                 string // job role
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
	GravatarURL     *string `json:"gravatar_url"`
	Position        *string `json:"position"`
}

// ValidatePassword accepts a potential password for a given user and attempts
// to validate it against the hash stored in the database after joining the
// supplied password with the stored password salt
func (u *User) ValidatePassword(password string) error {
	saltAndPass := []byte(fmt.Sprintf("%s%s", password, u.Salt))
	return bcrypt.CompareHashAndPassword(u.Password, saltAndPass)
}
