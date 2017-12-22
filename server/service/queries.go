package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) ListQueries(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Query, error) {
	return svc.ds.ListQueries(opt)
}

func (svc service) GetQuery(ctx context.Context, id uint) (*kolide.Query, error) {
	return svc.ds.Query(id)
}

func (svc service) NewQuery(ctx context.Context, p kolide.QueryPayload) (*kolide.Query, error) {
	query := &kolide.Query{Saved: true}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Description != nil {
		query.Description = *p.Description
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	vc, ok := viewer.FromContext(ctx)
	if ok {
		query.AuthorID = vc.UserID()
		query.AuthorName = vc.FullName()
	}

	query, err := svc.ds.NewQuery(query)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (svc service) ModifyQuery(ctx context.Context, id uint, p kolide.QueryPayload) (*kolide.Query, error) {
	query, err := svc.ds.Query(id)
	if err != nil {
		return nil, err
	}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Description != nil {
		query.Description = *p.Description
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	err = svc.ds.SaveQuery(query)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (svc service) DeleteQuery(ctx context.Context, id uint) error {
	return svc.ds.DeleteQuery(id)
}

func (svc service) DeleteQueries(ctx context.Context, ids []uint) (uint, error) {
	return svc.ds.DeleteQueries(ids)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeCreateQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeModifyQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeDeleteQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req deleteQueryRequest
	req.ID = id
	return req, nil
}

func decodeDeleteQueriesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req deleteQueriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getQueryRequest
	req.ID = id
	return req, nil
}

func decodeListQueriesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listQueriesRequest{ListOptions: opt}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type getQueryRequest struct {
	ID uint
}

type getQueryResponse struct {
	Query *kolide.Query `json:"query,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r getQueryResponse) error() error { return r.Err }

func makeGetQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getQueryRequest)
		query, err := svc.GetQuery(ctx, req.ID)
		if err != nil {
			return getQueryResponse{Err: err}, nil
		}
		return getQueryResponse{query, nil}, nil
	}
}

type listQueriesRequest struct {
	ListOptions kolide.ListOptions
}

type listQueriesResponse struct {
	Queries []kolide.Query `json:"queries"`
	Err     error          `json:"error,omitempty"`
}

func (r listQueriesResponse) error() error { return r.Err }

func makeListQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listQueriesRequest)
		queries, err := svc.ListQueries(ctx, req.ListOptions)
		if err != nil {
			return listQueriesResponse{Err: err}, nil
		}

		resp := listQueriesResponse{Queries: []kolide.Query{}}
		for _, query := range queries {
			resp.Queries = append(resp.Queries, *query)
		}
		return resp, nil
	}
}

type createQueryRequest struct {
	payload kolide.QueryPayload
}

type createQueryResponse struct {
	Query *kolide.Query `json:"query,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r createQueryResponse) error() error { return r.Err }

func makeCreateQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createQueryRequest)
		query, err := svc.NewQuery(ctx, req.payload)
		if err != nil {
			return createQueryResponse{Err: err}, nil
		}
		return createQueryResponse{query, nil}, nil
	}
}

type modifyQueryRequest struct {
	ID      uint
	payload kolide.QueryPayload
}

type modifyQueryResponse struct {
	Query *kolide.Query `json:"query,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r modifyQueryResponse) error() error { return r.Err }

func makeModifyQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyQueryRequest)
		query, err := svc.ModifyQuery(ctx, req.ID, req.payload)
		if err != nil {
			return modifyQueryResponse{Err: err}, nil
		}
		return modifyQueryResponse{query, nil}, nil
	}
}

type deleteQueryRequest struct {
	ID uint
}

type deleteQueryResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteQueryResponse) error() error { return r.Err }

func makeDeleteQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteQueryRequest)
		err := svc.DeleteQuery(ctx, req.ID)
		if err != nil {
			return deleteQueryResponse{Err: err}, nil
		}
		return deleteQueryResponse{}, nil
	}
}

type deleteQueriesRequest struct {
	IDs []uint `json:"ids"`
}

type deleteQueriesResponse struct {
	Deleted uint  `json:"deleted"`
	Err     error `json:"error,omitempty"`
}

func (r deleteQueriesResponse) error() error { return r.Err }

func makeDeleteQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteQueriesRequest)
		deleted, err := svc.DeleteQueries(ctx, req.IDs)
		if err != nil {
			return deleteQueriesResponse{Err: err}, nil
		}
		return deleteQueriesResponse{Deleted: deleted}, nil
	}
}
