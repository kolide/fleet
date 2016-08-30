package kitserver

import (
	"net/http"

	"golang.org/x/net/context"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/kolide/kolide-ose/kolide"
)

// MakeHandler creates an http handler for the Kolide API
func MakeHandler(ctx context.Context, svc kolide.Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(
			setViewerContext(svc, logger),
		),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerAfter(
			kithttp.SetContentType("application/json; charset=utf-8"),
		),
	}

	createUserHandler := kithttp.NewServer(
		ctx,
		makeCreateUserEndpoint(svc),
		decodeCreateUserRequest,
		encodeResponse,
		opts...,
	)

	getUserHandler := kithttp.NewServer(
		ctx,
		makeGetUserEndpoint(svc),
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)

	changePasswordHandler := kithttp.NewServer(
		ctx,
		makeChangePasswordEndpoint(svc),
		decodeChangePasswordRequest,
		encodeResponse,
		opts...,
	)

	updateAdminRoleHandler := kithttp.NewServer(
		ctx,
		makeUpdateAdminRoleEndpoint(svc),
		decodeUpdateAdminRoleRequest,
		encodeResponse,
		opts...,
	)

	updateUserStatusHandler := kithttp.NewServer(
		ctx,
		makeUpdateUserStatusEndpoint(svc),
		decodeUpdateUserStatusRequest,
		encodeResponse,
		opts...,
	)

	api := mux.NewRouter()
	api.Handle("/api/v1/kolide/users", createUserHandler).Methods("POST")
	api.Handle("/api/v1/kolide/users/{id}", getUserHandler).Methods("GET")
	api.Handle("/api/v1/kolide/users/{id}/password", changePasswordHandler).Methods("POST")
	api.Handle("/api/v1/kolide/users/{id}/role", updateAdminRoleHandler).Methods("POST")
	api.Handle("/api/v1/kolide/users/{id}/status", updateUserStatusHandler).Methods("POST")

	r := mux.NewRouter()

	r.PathPrefix("/api/v1/kolide").Handler(authMiddleware(svc, logger, api))
	r.Handle("/login", login(svc, logger)).Methods("POST")
	r.Handle("/logout", logout(svc, logger)).Methods("GET")
	return r
}

// setViewerContext updates the context with a viewerContext,
// which holds the currently logged in user
func setViewerContext(svc kolide.Service, logger kitlog.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		sm := svc.NewSessionManager(ctx, nil, r)
		session, err := sm.Session()
		if err != nil {
			logger.Log("err", err, "error-source", "setViewerContext")
			return ctx
		}

		user, err := svc.User(ctx, session.UserID)
		if err != nil {
			logger.Log("err", err, "error-source", "setViewerContext")
			return ctx
		}

		ctx = context.WithValue(ctx, "viewerContext", &viewerContext{
			user: user,
		})
		return ctx
	}
}
