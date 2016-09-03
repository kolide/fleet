package kitserver

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Create Pack
////////////////////////////////////////////////////////////////////////////////

type createPackRequest struct {
	payload kolide.PackPayload
}

type createPackResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Err      error  `json:"error, omitempty"`
}

func (r createPackResponse) error() error { return r.Err }

func makeCreatePackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createPackRequest)
		pack, err := svc.NewPack(ctx, req.payload)
		if err != nil {
			return createPackResponse{Err: err}, nil
		}
		return createPackResponse{
			ID:       pack.ID,
			Name:     pack.Name,
			Platform: pack.Platform,
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify Pack
////////////////////////////////////////////////////////////////////////////////

type modifyPackRequest struct {
	ID      uint
	payload kolide.PackPayload
}

type modifyPackResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Err      error  `json:"error, omitempty"`
}

func (r modifyPackResponse) error() error { return r.Err }

func makeModifyPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyPackRequest)
		pack, err := svc.ModifyPack(ctx, req.ID, req.payload)
		if err != nil {
			return modifyPackResponse{Err: err}, nil
		}
		return modifyPackResponse{
			ID:       pack.ID,
			Name:     pack.Name,
			Platform: pack.Platform,
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Pack
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Get Pack
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Get All Packs
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Add Query To Pack
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Delete Query From Pack
////////////////////////////////////////////////////////////////////////////////
