package kitserver

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/net/context"
)

type createQueryRequest struct {
	payload kolide.QueryPayload
}

type createQueryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Query        string `json:"query"`
	Interval     uint   `json:"interval"`
	Snapshot     bool   `json:"snapshot"`
	Differential bool   `json:"differential"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
	Err          error  `json:"error, omitempty"`
}

func (r createQueryResponse) error() error { return r.Err }

func makeCreateQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createQueryRequest)
		query, err := svc.NewQuery(ctx, req.payload)
		if err != nil {
			return createQueryResponse{Err: err}, nil
		}
		return createQueryResponse{
			ID:           query.ID,
			Name:         query.Name,
			Query:        query.Query,
			Interval:     query.Interval,
			Snapshot:     query.Snapshot,
			Differential: query.Differential,
			Platform:     query.Platform,
			Version:      query.Version,
		}, nil
	}
}

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
