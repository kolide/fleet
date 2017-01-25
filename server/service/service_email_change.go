package service

import "golang.org/x/net/context"

func (svc service) CommitEmailChange(ctx context.Context, token string) (string, error) {
	return svc.ds.CommitEmailChange(token)
}
