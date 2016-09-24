package mock

import (
	"golang.org/x/net/context"

	"github.com/kolide/kolide-ose/kolide"
)

var _ kolide.UserService = (*UserService)(nil)

// UserService is a mock struct which implements
// kolide.UserService for use in tests
type UserService struct {
	NewUserFunc        func(ctx context.Context, p kolide.UserPayload) (*kolide.User, error)
	NewUserFuncInvoked bool

	UserFunc        func(ctx context.Context, id uint) (*kolide.User, error)
	UserFuncInvoked bool

	UsersFunc        func(ctx context.Context) ([]*kolide.User, error)
	UsersFuncInvoked bool

	ModifyUserFunc        func(ctx context.Context, userID uint, p kolide.UserPayload) (*kolide.User, error)
	ModifyUserFuncInvoked bool

	AuthenticatedUserFunc        func(ctx context.Context) (*kolide.User, error)
	AuthenticatedUserFuncInvoked bool

	RequestPasswordResetFunc        func(ctx context.Context, email string) error
	RequestPasswordResetFuncInvoked bool

	ResetPasswordFunc        func(ctx context.Context, token string, password string) error
	ResetPasswordFuncInvoked bool
}

func (svc *UserService) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	svc.NewUserFuncInvoked = true
	return svc.NewUserFunc(ctx, p)
}

func (svc *UserService) User(ctx context.Context, id uint) (*kolide.User, error) {
	svc.UserFuncInvoked = true
	return svc.UserFunc(ctx, id)
}

func (svc *UserService) Users(ctx context.Context) ([]*kolide.User, error) {
	svc.UsersFuncInvoked = true
	return svc.UsersFunc(ctx)
}

func (svc *UserService) ModifyUser(ctx context.Context, userID uint, p kolide.UserPayload) (*kolide.User, error) {
	svc.ModifyUserFuncInvoked = true
	return svc.ModifyUserFunc(ctx, userID, p)
}

func (svc *UserService) AuthenticatedUser(ctx context.Context) (*kolide.User, error) {
	svc.AuthenticatedUserFuncInvoked = true
	return svc.AuthenticatedUserFunc(ctx)
}

func (svc *UserService) RequestPasswordReset(ctx context.Context, email string) error {
	svc.RequestPasswordResetFuncInvoked = true
	return svc.RequestPasswordResetFunc(ctx, email)
}

func (svc *UserService) ResetPassword(ctx context.Context, token string, password string) error {
	svc.ResetPasswordFuncInvoked = true
	return svc.ResetPasswordFunc(ctx, token, password)
}
