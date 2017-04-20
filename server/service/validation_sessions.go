package service

import "context"

func (mw validationMiddleware) CallbackSSO(ctx context.Context, ssoHandle, userID string) (string, error) {
	invalid := &invalidArgumentError{}
	if userID == "" {
		invalid.Append("SAMLResponse", "missing user ID")
	}
	if ssoHandle == "" {
		invalid.Append("RelayState", "missing required relay state")
	}
	if invalid.HasErrors() {
		return "", invalid
	}
	return mw.Service.CallbackSSO(ctx, ssoHandle, userID)
}
