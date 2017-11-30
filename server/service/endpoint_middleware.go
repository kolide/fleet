package service

import (
	"context"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/contexts/token"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
)

var errNoContext = errors.New("context key not set")

// authenticatedUser wraps an endpoint, requires that the Kolide user is
// authenticated, and populates the context with a Viewer struct for that user.
func authenticatedUser(jwtKey string, svc kolide.Service, next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// first check if already successfully set
		if _, ok := viewer.FromContext(ctx); ok {
			return next(ctx, request)
		}

		// if not succesful, try again this time with errors
		bearer, ok := token.FromContext(ctx)
		if !ok {
			return nil, authError{reason: "no auth token"}
		}

		v, err := authViewer(ctx, jwtKey, bearer, svc)
		if err != nil {
			return nil, err
		}

		ctx = viewer.NewContext(ctx, *v)
		return next(ctx, request)
	}
}

// authViewer creates an authenticated viewer by validating a JWT token.
func authViewer(ctx context.Context, jwtKey string, bearerToken token.Token, svc kolide.Service) (*viewer.Viewer, error) {
	jwtToken, err := jwt.Parse(string(bearerToken), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, authError{reason: err.Error()}
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, authError{reason: "no jwt claims"}
	}
	sessionKeyClaim, ok := claims["session_key"]
	if !ok {
		return nil, authError{reason: "no session_key in JWT claims"}
	}
	sessionKey, ok := sessionKeyClaim.(string)
	if !ok {
		return nil, authError{reason: "non-string key in sessionClaim"}
	}
	session, err := svc.GetSessionByKey(ctx, sessionKey)
	if err != nil {
		return nil, authError{reason: err.Error()}
	}
	user, err := svc.User(ctx, session.UserID)
	if err != nil {
		return nil, authError{reason: err.Error()}
	}
	return &viewer.Viewer{User: user, Session: session}, nil
}

func mustBeAdmin(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, errNoContext
		}
		if !vc.IsAdmin() {
			return nil, permissionError{message: "must be an admin"}
		}
		return next(ctx, request)
	}
}

func canPerformActions(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, errNoContext
		}
		if !vc.CanPerformActions() {
			return nil, permissionError{message: "no read permissions"}
		}
		return next(ctx, request)
	}
}

func canReadUser(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, errNoContext
		}
		uid := requestUserIDFromContext(ctx)
		if !vc.CanPerformReadActionOnUser(uid) {
			return nil, permissionError{message: "no read permissions on user"}
		}
		return next(ctx, request)
	}
}

func canModifyUser(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, errNoContext
		}
		uid := requestUserIDFromContext(ctx)
		if !vc.CanPerformWriteActionOnUser(uid) {
			return nil, permissionError{message: "no write permissions on user"}
		}
		return next(ctx, request)
	}
}

type permission int

const (
	anyone permission = iota
	self
	admin
)

func requestUserIDFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value("request-id").(uint)
	if !ok {
		return 0
	}
	return userID
}
