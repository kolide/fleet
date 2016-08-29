package kitserver

import (
	"golang.org/x/net/context"

	"github.com/kolide/kolide-ose/kolide"
)

type validationMiddleware struct {
	kolide.Service
}

func (mw validationMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	// check required params
	if p.Username == nil {
		return nil, invalidArgumentError{field: "username"}
	}

	if p.Password == nil {
		return nil, invalidArgumentError{field: "password"}
	}

	if p.Email == nil {
		return nil, invalidArgumentError{field: "email"}
	}

	return mw.Service.NewUser(ctx, p)
}
