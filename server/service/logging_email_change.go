package service

import (
	"time"

	"golang.org/x/net/context"
)

func (mw loggingMiddleware) CommitEmailChange(ctx context.Context, token string) (string, error) {
	var (
		err     error
		newMail string
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method",
			"CommitEmailChange",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	newMail, err = mw.Service.CommitEmailChange(ctx, token)
	return newMail, err
}
