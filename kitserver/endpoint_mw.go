package kitserver

import (
	"errors"

	"github.com/go-kit/kit/endpoint"

	"golang.org/x/net/context"
)

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

func viewerFromContext(ctx context.Context) (*ViewerContext, error) {
	vc, ok := ctx.Value("viewerContext").(*ViewerContext)
	if !ok {
		return nil, errors.New("no viewer context set")
	}
	return vc, nil
}
