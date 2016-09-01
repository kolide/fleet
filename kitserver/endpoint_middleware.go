package kitserver

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

var errNoContext = errors.New("no viewer context set")

func mustBeAdmin(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, err := viewerFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if !vc.IsAdmin() {
			return nil, forbiddenError{message: "must be an admin"}
		}
		return next(ctx, request)
	}
}

func canReadUser(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, err := viewerFromContext(ctx)
		if err != nil {
			return nil, err
		}
		uid := requestUserIDFromContext(ctx)
		// TODO discuss the semantics of this check
		if !vc.CanPerformReadActionOnUser(uid) {
			return nil, forbiddenError{message: "no read permissions on user"}
		}
		return next(ctx, request)
	}
}

func canModifyUser(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, err := viewerFromContext(ctx)
		if err != nil {
			return nil, err
		}
		uid := requestUserIDFromContext(ctx)
		if !vc.CanPerformWriteActionOnUser(uid) {
			return nil, forbiddenError{message: "no write permissions on user"}
		}
		return next(ctx, request)
	}
}

func requestUserIDFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value("request-id").(uint)
	if !ok {
		return 0
	}
	return userID
}

func viewerFromContext(ctx context.Context) (*ViewerContext, error) {
	vc, ok := ctx.Value("viewerContext").(*ViewerContext)
	if !ok {
		return nil, errNoContext
	}
	return vc, nil
}
