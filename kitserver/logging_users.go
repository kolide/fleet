package kitserver

import (
	"time"

	"golang.org/x/net/context"

	"github.com/kolide/kolide-ose/kolide"
)

func (mw loggingMiddleware) NewUser(ctx context.Context, p kolide.UserPayload) (user *kolide.User, err error) {
	vc, err := viewerFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var username = "none"

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "NewUser",
			"user", username,
			"err", err,
			"created_by", vc.user.Username,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.NewUser(ctx, p)

	if user != nil {
		username = user.Username
	}
	return
}

func (mw loggingMiddleware) User(ctx context.Context, id uint) (user *kolide.User, err error) {
	var username = "none"

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "User",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, err = mw.Service.User(ctx, id)

	if user != nil {
		username = user.Username
	}
	return
}
