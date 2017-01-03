package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type deleteEntityRequest struct {
	Entity kolide.Entity
	ID     uint
}

func (d *deleteEntityRequest) EntityType() string {
	return kolide.DBTable(d.Entity)
}

func (d *deleteEntityRequest) EntityID() uint {
	return d.ID
}

type deleteEntityResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteEntityResponse) error() error { return r.Err }

func makeDeleteEntityEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteEntityRequest)
		err := svc.Delete(ctx, &req)
		if err != nil {
			return deleteEntityResponse{Err: err}, nil
		}
		return deleteEntityResponse{}, nil
	}
}
