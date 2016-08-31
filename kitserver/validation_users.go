package kitserver

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/kolide/kolide-ose/kolide"
)

type validationMiddleware struct {
	kolide.Service
}

func (mw validationMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (*kolide.User, error) {
	if err := mw.authCheckAdmin(ctx); err != nil {
		return nil, err
	}
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

func (mw validationMiddleware) authCheckAdmin(ctx context.Context) error {
	vc, ok := ctx.Value("viewerContext").(viewerContext)
	if !ok {
		return errors.New("no viewer context set")
	}
	if !vc.IsAdmin() {
		return forbiddenError{message: "must be an admin"}
	}
	return nil
}
