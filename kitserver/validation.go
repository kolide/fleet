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
