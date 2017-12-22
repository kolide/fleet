package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeChangeEmailRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		return nil, errBadRoute
	}

	response := changeEmailRequest{
		Token: token,
	}

	return response, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type changeEmailRequest struct {
	Token string
}

type changeEmailResponse struct {
	NewEmail string `json:"new_email"`
	Err      error  `json:"error,omitempty"`
}

func (r changeEmailResponse) error() error { return r.Err }

func makeChangeEmailEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(changeEmailRequest)
		newEmailAddress, err := svc.ChangeUserEmail(ctx, req.Token)
		if err != nil {
			return changeEmailResponse{Err: err}, nil
		}
		return changeEmailResponse{NewEmail: newEmailAddress}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) ChangeUserEmail(ctx context.Context, token string) (string, error) {
	var (
		err      error
		newEmail string
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "CommitEmailChange", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	newEmail, err = mw.Service.ChangeUserEmail(ctx, token)
	return newEmail, err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) ChangeUserEmail(ctx context.Context, token string) (string, error) {
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
	newMail, err = mw.Service.ChangeUserEmail(ctx, token)
	return newMail, err
}
