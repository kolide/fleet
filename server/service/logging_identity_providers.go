package service

import (
	"context"
	"time"

	"github.com/kolide/kolide/server/kolide"
)

func (mw loggingMiddleware) GetIdentityProvider(ctx context.Context, id uint) (*kolide.IdentityProvider, error) {
	var (
		idp *kolide.IdentityProvider
		err error
	)

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetIdentityProvider",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	idp, err = mw.Service.GetIdentityProvider(ctx, id)
	return idp, err
}

func (mw loggingMiddleware) ModifyIdentityProvider(ctx context.Context, id uint, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var (
		idp *kolide.IdentityProvider
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ModifyIdentityProvider",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	idp, err = mw.Service.ModifyIdentityProvider(ctx, id, payload)
	return idp, err
}

func (mw loggingMiddleware) NewIdentityProvider(ctx context.Context, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var (
		idp *kolide.IdentityProvider
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "NewIdentityProvider",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	idp, err = mw.Service.NewIdentityProvider(ctx, payload)
	return idp, err
}

func (mw loggingMiddleware) DeleteIdentityProvider(ctx context.Context, id uint) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "DeleteIdentityProvider",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Service.DeleteIdentityProvider(ctx, id)
	return err
}

func (mw loggingMiddleware) ListIdentityProviders(ctx context.Context) ([]kolide.IdentityProvider, error) {
	var (
		idps []kolide.IdentityProvider
		err  error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ListIdentityProviders",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	idps, err = mw.Service.ListIdentityProviders(ctx)
	return idps, err
}
