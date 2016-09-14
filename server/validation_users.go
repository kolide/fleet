package server

import (
	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/net/context"
)

type validationMiddleware struct {
	kolide.Service
}

func (mw validationMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	// check required params
	if p.Username == nil {
		return nil, invalidArgumentError{field: "username", required: true}
	}

	if p.Password == nil {
		return nil, invalidArgumentError{field: "password", required: true}
	}

	if p.Email == nil {
		return nil, invalidArgumentError{field: "email", required: true}
	}

	return mw.Service.NewUser(ctx, p)
}

func (mw validationMiddleware) ChangePassword(ctx context.Context, userID uint, token, password string) error {

	vc, err := viewerContextFromContext(ctx)
	if err != nil {
		return err
	}

	// require token unless an admin forced the password reset
	if token == "" && vc.user == nil {
		return invalidArgumentError{field: "token", required: true}
	}

	if password == "" {
		return invalidArgumentError{field: "password", required: true}
	}
	return mw.Service.ChangePassword(ctx, userID, token, password)
}
