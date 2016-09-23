package mock

import (
	"golang.org/x/net/context"

	"github.com/kolide/kolide-ose/kolide"
)

var _ kolide.SessionService = (*SessionService)(nil)

type SessionService struct {
	LoginFunc        func(ctx context.Context, username string, password string) (*kolide.User, string, error)
	LoginFuncInvoked bool

	LogoutFunc        func(ctx context.Context) error
	LogoutFuncInvoked bool

	DestroySessionFunc        func(ctx context.Context) error
	DestroySessionFuncInvoked bool

	GetInfoAboutSessionsForUserFunc        func(ctx context.Context, id uint) ([]*kolide.Session, error)
	GetInfoAboutSessionsForUserFuncInvoked bool

	DeleteSessionsForUserFunc        func(ctx context.Context, id uint) error
	DeleteSessionsForUserFuncInvoked bool

	GetInfoAboutSessionFunc        func(ctx context.Context, id uint) (*kolide.Session, error)
	GetInfoAboutSessionFuncInvoked bool

	GetSessionByKeyFunc        func(ctx context.Context, key string) (*kolide.Session, error)
	GetSessionByKeyFuncInvoked bool

	DeleteSessionFunc        func(ctx context.Context, id uint) error
	DeleteSessionFuncInvoked bool
}

func (svc *SessionService) Login(ctx context.Context, username string, password string) (*kolide.User, string, error) {
	svc.LoginFuncInvoked = true
	return svc.LoginFunc(ctx, username, password)
}

func (svc *SessionService) Logout(ctx context.Context) error {
	svc.LogoutFuncInvoked = true
	return svc.LogoutFunc(ctx)
}

func (svc *SessionService) DestroySession(ctx context.Context) error {
	svc.DestroySessionFuncInvoked = true
	return svc.DestroySessionFunc(ctx)
}

func (svc *SessionService) GetInfoAboutSessionsForUser(ctx context.Context, id uint) ([]*kolide.Session, error) {
	svc.GetInfoAboutSessionsForUserFuncInvoked = true
	return svc.GetInfoAboutSessionsForUserFunc(ctx, id)
}

func (svc *SessionService) DeleteSessionsForUser(ctx context.Context, id uint) error {
	svc.DeleteSessionsForUserFuncInvoked = true
	return svc.DeleteSessionsForUserFunc(ctx, id)
}

func (svc *SessionService) GetInfoAboutSession(ctx context.Context, id uint) (*kolide.Session, error) {
	svc.GetInfoAboutSessionFuncInvoked = true
	return svc.GetInfoAboutSessionFunc(ctx, id)
}

func (svc *SessionService) GetSessionByKey(ctx context.Context, key string) (*kolide.Session, error) {
	svc.GetSessionByKeyFuncInvoked = true
	return svc.GetSessionByKeyFunc(ctx, key)
}

func (svc *SessionService) DeleteSession(ctx context.Context, id uint) error {
	svc.DeleteSessionFuncInvoked = true
	return svc.DeleteSessionFunc(ctx, id)
}
