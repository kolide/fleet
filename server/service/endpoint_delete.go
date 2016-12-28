package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type deleteEntityRequest struct {
	EntityType string
	ID         uint
}

type deleteEntityResponse struct {
	Err error `json:"error,omitempty"`
}

func makeDeleteEntityEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteEntityRequest)
		e := &entity{req.EntityType, req.ID}
		err := svc.Delete(ctx, e)
		if err != nil {
			return deleteEntityResponse{Err: err}, nil
		}
		return deleteEntityResponse{}, nil
	}
}
