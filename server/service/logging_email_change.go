package service

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"time"
)

func (mw loggingMiddleware) ChangeUserEmail(ctx context.Context, token string) (string, error) {
	var (
		err     error
		newMail string
	)
	defer func(begin time.Time) {
		_ = level.Debug(mw.logger).Log(
			"method",
			"CommitEmailChange",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	newMail, err = mw.Service.ChangeUserEmail(ctx, token)
	return newMail, err
}
