package service

import (
	"context"

	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/sso"
	"github.com/pkg/errors"
)

func (mw validationMiddleware) CallbackSSO(ctx context.Context, auth kolide.Auth) (string, error) {
	invalid := &invalidArgumentError{}
	status, err := auth.Status()
	if err != nil {
		return "", errors.Wrap(err, "CallbackSSO validation")
	}
	if status == sso.Success {
		if auth.UserID() == "" {
			invalid.Append("SAMLResponse", "missing user ID")
		}
		if auth.RelayState() == "" {
			invalid.Append("RelayState", "missing required relay state")
		}
	}
	if invalid.HasErrors() {
		return "", invalid
	}
	return mw.Service.CallbackSSO(ctx, auth)
}
