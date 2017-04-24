package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kolide/kolide/server/kolide"
)

func (mw metricsMiddleware) NewIdentityProvider(ctx context.Context, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var (
		idp *kolide.IdentityProvider
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "NewIdentityProvider", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	idp, err = mw.Service.NewIdentityProvider(ctx, payload)
	return idp, err
}

func (mw metricsMiddleware) ModifyIdentityProvider(ctx context.Context, id uint, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var (
		idp *kolide.IdentityProvider
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyIdentityProvider", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	idp, err = mw.Service.ModifyIdentityProvider(ctx, id, payload)
	return idp, err
}

func (mw metricsMiddleware) GetIdentityProvider(ctx context.Context, id uint) (*kolide.IdentityProvider, error) {
	var (
		idp *kolide.IdentityProvider
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "GetIdentityProvider", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	idp, err = mw.Service.GetIdentityProvider(ctx, id)
	return idp, err
}

func (mw metricsMiddleware) ListIdentityProviders(ctx context.Context) ([]kolide.IdentityProvider, error) {
	var (
		idps []kolide.IdentityProvider
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ListIdentityProviders", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	idps, err = mw.Service.ListIdentityProviders(ctx)
	return idps, err
}

func (mw metricsMiddleware) ListIdentityProvidersNoAuth(ctx context.Context) ([]kolide.IdentityProviderNoAuth, error) {
	var (
		idps []kolide.IdentityProviderNoAuth
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ListIdentityProvidersNoAuth", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	idps, err = mw.Service.ListIdentityProvidersNoAuth(ctx)
	return idps, err
}

func (mw metricsMiddleware) DeleteIdentityProvider(ctx context.Context, id uint) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteIdentityProvider", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.DeleteIdentityProvider(ctx, id)
	return err
}
