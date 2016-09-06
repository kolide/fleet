package server

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/net/context"
)

type enrollOsqueryRequest struct {
	EnrollSecret   string `json:"secret"`
	HostIdentifier string `json:"host_identifier"`
}

type enrollOsqueryResponse struct {
	NodeKey     string `json:"node_key,omitempty"`
	NodeInvalid bool   `json:"node_invalid,omitempty"`
	Err         error  `json:"error,omitempty"`
}

func (r enrollOsqueryResponse) error() error { return r.Err }

func makeEnrollOsqueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(enrollOsqueryRequest)
		nodeKey, err := svc.EnrollAgent(ctx, req.EnrollSecret, req.HostIdentifier)
		if err != nil {
			return enrollOsqueryResponse{
				NodeInvalid: true,
				Err:         err,
			}, nil
		}
		return enrollOsqueryResponse{
			NodeKey: nodeKey,
		}, nil
	}
}
