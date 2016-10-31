package service

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// List Targrets
////////////////////////////////////////////////////////////////////////////////

type listTargetsRequest struct {
	Query string `json:"query"`
}

type targetsData struct {
	Hosts  []hostResponse `json:"hosts"`
	Labels []kolide.Label `json:"labels"`
}

type listTargetsResponse struct {
	Targets *targetsData `json:"targets,omitempty"`
	Err     error        `json:"error,omitempty"`
}

func (r listTargetsResponse) error() error { return r.Err }

func makeListTargetsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, errors.New("Unimplemented")
	}
}
