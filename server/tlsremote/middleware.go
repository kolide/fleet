package tlsremote

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/http"

	"github.com/kolide/fleet/server/kolide"
)

type Middleware func(kolide.OsqueryService) kolide.OsqueryService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next kolide.OsqueryService) kolide.OsqueryService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   kolide.OsqueryService
	logger log.Logger
}

func (mw loggingMiddleware) EnrollAgent(ctx context.Context, enrollSecret string, hostIdentifier string) (string, error) {
	var (
		nodeKey string
		err     error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "EnrollAgent",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	nodeKey, err = mw.next.EnrollAgent(ctx, enrollSecret, hostIdentifier)
	return nodeKey, err
}

func (mw loggingMiddleware) AuthenticateHost(ctx context.Context, nodeKey string) (*kolide.Host, error) {
	var (
		host *kolide.Host
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "AuthenticateHost",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	host, err = mw.next.AuthenticateHost(ctx, nodeKey)
	return host, err
}

func (mw loggingMiddleware) GetClientConfig(ctx context.Context) (*kolide.OsqueryConfig, error) {
	var (
		config *kolide.OsqueryConfig
		err    error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetClientConfig",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	config, err = mw.next.GetClientConfig(ctx)
	return config, err
}

func (mw loggingMiddleware) GetDistributedQueries(ctx context.Context) (map[string]string, uint, error) {
	var (
		queries    map[string]string
		err        error
		accelerate uint
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetDistributedQueries",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	queries, accelerate, err = mw.next.GetDistributedQueries(ctx)
	return queries, accelerate, err
}

func (mw loggingMiddleware) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]string) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitDistributedQueryResults",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.SubmitDistributedQueryResults(ctx, results, statuses)
	return err
}

func (mw loggingMiddleware) SubmitStatusLogs(ctx context.Context, logs []kolide.OsqueryStatusLog) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitStatusLogs",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.SubmitStatusLogs(ctx, logs)
	return err
}

func (mw loggingMiddleware) SubmitResultLogs(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitResultLogs",
			"ip_addr", ctx.Value(http.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.SubmitResultLogs(ctx, logs)
	return err
}
