package service

import (
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func decodeDeleteEntityRequest(entityType string) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		id, err := idFromRequest(r, "id")
		if err != nil {
			return nil, err
		}
		return deleteEntityRequest{EntityType: entityType, ID: id}, nil
	}
}
