package kolide

import "time"

// PasswordResetStore manages password resets in the Datastore
type PasswordResetStore interface {
	NewPasswordResetRequest(req *PasswordResetRequest) (*PasswordResetRequest, error)
	SavePasswordResetRequest(req *PasswordResetRequest) error
	DeletePasswordResetRequest(req *PasswordResetRequest) error
	DeletePasswordResetRequestsForUser(userID uint) error
	FindPassswordResetByID(id uint) (*PasswordResetRequest, error)
	FindPassswordResetsByUserID(id uint) ([]*PasswordResetRequest, error)
	FindPassswordResetByToken(token string) (*PasswordResetRequest, error)
	FindPassswordResetByTokenAndUserID(token string, id uint) (*PasswordResetRequest, error)
}

// Campaign is an email campaign
// Types which implement the Campaign interface
// can be marshalled into an email body
type Campaign interface {
	Message() ([]byte, error)
}

type Email struct {
	To   []string
	From string
	msg  Campaign
}

type MailService interface {
	SendEmail(e Email) error
}

// PasswordResetRequest represents a database table for
// Password Reset Requests
type PasswordResetRequest struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	UserID    uint
	Token     string `gorm:"size:1024"`
}

func (r PasswordResetRequest) Message() ([]byte, error) {
	// TODO: marshal error into the correct body
	msg := []byte("temporary")
	return msg, nil
}
