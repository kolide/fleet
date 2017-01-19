package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type certificateRequest struct {
	Insecure bool
}

type certificateResponse struct {
	CertificateChain []byte `json:"certificate_chain"`
	Err              error  `json:"error,omitempty"`
}

func (r certificateResponse) error() error { return r.Err }

func makeCertificateEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(certificateRequest)
		chain, err := svc.CertificateChain(ctx, req.Insecure)
		if err != nil {
			return certificateResponse{Err: err}, nil
		}
		return certificateResponse{CertificateChain: chain}, nil
	}
}
