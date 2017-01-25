package service

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

func (mw metricsMiddleware) CommitEmailChange(ctx context.Context, token string) (string, error) {
	var (
		err      error
		newEmail string
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "CommitEmailChange", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	newEmail, err = mw.Service.CommitEmailChange(ctx, token)
	return newEmail, err
}
