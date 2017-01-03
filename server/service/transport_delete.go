package service

import (
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func decodeDeleteEntityRequest(e kolide.Entity) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		id, err := idFromRequest(r, "id")
		if err != nil {
			return nil, err
		}
		return deleteEntityRequest{Entity: e, ID: id}, nil
	}
}
