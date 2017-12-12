package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) GetScheduledQuery(ctx context.Context, id uint) (*kolide.ScheduledQuery, error) {
	return svc.ds.ScheduledQuery(id)
}

func (svc service) GetScheduledQueriesInPack(ctx context.Context, id uint, opts kolide.ListOptions) ([]*kolide.ScheduledQuery, error) {
	return svc.ds.ListScheduledQueriesInPack(id, opts)
}

func (svc service) ScheduleQuery(ctx context.Context, sq *kolide.ScheduledQuery) (*kolide.ScheduledQuery, error) {
	return svc.ds.NewScheduledQuery(sq)
}

func (svc service) ModifyScheduledQuery(ctx context.Context, id uint, p kolide.ScheduledQueryPayload) (*kolide.ScheduledQuery, error) {
	sq, err := svc.GetScheduledQuery(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "getting scheduled query to modify")
	}

	if p.PackID != nil {
		sq.PackID = *p.PackID
	}

	if p.QueryID != nil {
		sq.QueryID = *p.QueryID
	}

	if p.Interval != nil {
		sq.Interval = *p.Interval
	}

	if p.Snapshot != nil {
		sq.Snapshot = p.Snapshot
	}

	if p.Removed != nil {
		sq.Removed = p.Removed
	}

	if p.Platform != nil {
		sq.Platform = p.Platform
	}

	if p.Version != nil {
		sq.Version = p.Version
	}

	if p.Shard != nil {
		sq.Shard = p.Shard
	}

	return svc.ds.SaveScheduledQuery(sq)
}

func (svc service) DeleteScheduledQuery(ctx context.Context, id uint) error {
	return svc.ds.DeleteScheduledQuery(id)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeScheduleQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req scheduleQueryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeModifyScheduledQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyScheduledQueryRequest

	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	req.ID = id
	return req, nil
}

func decodeDeleteScheduledQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req deleteScheduledQueryRequest
	req.ID = id
	return req, nil
}

func decodeGetScheduledQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getScheduledQueryRequest
	req.ID = id
	return req, nil
}

func decodeGetScheduledQueriesInPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getScheduledQueriesInPackRequest
	req.ID = id
	return req, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type getScheduledQueryRequest struct {
	ID uint
}

type scheduledQueryResponse struct {
	kolide.ScheduledQuery
}

type getScheduledQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r getScheduledQueryResponse) error() error { return r.Err }

func makeGetScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getScheduledQueryRequest)

		sq, err := svc.GetScheduledQuery(ctx, req.ID)
		if err != nil {
			return getScheduledQueryResponse{Err: err}, nil
		}

		return getScheduledQueryResponse{
			Scheduled: scheduledQueryResponse{
				ScheduledQuery: *sq,
			},
		}, nil
	}
}

type getScheduledQueriesInPackRequest struct {
	ID          uint
	ListOptions kolide.ListOptions
}

type getScheduledQueriesInPackResponse struct {
	Scheduled []scheduledQueryResponse `json:"scheduled"`
	Err       error                    `json:"error,omitempty"`
}

func (r getScheduledQueriesInPackResponse) error() error { return r.Err }

func makeGetScheduledQueriesInPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getScheduledQueriesInPackRequest)
		resp := getScheduledQueriesInPackResponse{Scheduled: []scheduledQueryResponse{}}

		queries, err := svc.GetScheduledQueriesInPack(ctx, req.ID, req.ListOptions)
		if err != nil {
			return getScheduledQueriesInPackResponse{Err: err}, nil
		}

		for _, q := range queries {
			resp.Scheduled = append(resp.Scheduled, scheduledQueryResponse{
				ScheduledQuery: *q,
			})
		}

		return resp, nil
	}
}

type scheduleQueryRequest struct {
	PackID   uint    `json:"pack_id"`
	QueryID  uint    `json:"query_id"`
	Interval uint    `json:"interval"`
	Snapshot *bool   `json:"snapshot"`
	Removed  *bool   `json:"removed"`
	Platform *string `json:"platform"`
	Version  *string `json:"version"`
	Shard    *uint   `json:"shard"`
}

type scheduleQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled"`
	Err       error                  `json:"error,omitempty"`
}

func makeScheduleQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(scheduleQueryRequest)

		scheduled, err := svc.ScheduleQuery(ctx, &kolide.ScheduledQuery{
			PackID:   req.PackID,
			QueryID:  req.QueryID,
			Interval: req.Interval,
			Snapshot: req.Snapshot,
			Removed:  req.Removed,
			Platform: req.Platform,
			Version:  req.Version,
			Shard:    req.Shard,
		})
		if err != nil {
			return scheduleQueryResponse{Err: err}, nil
		}
		return scheduleQueryResponse{Scheduled: scheduledQueryResponse{
			ScheduledQuery: *scheduled,
		}}, nil
	}
}

type modifyScheduledQueryRequest struct {
	ID      uint
	payload kolide.ScheduledQueryPayload
}

type modifyScheduledQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r modifyScheduledQueryResponse) error() error { return r.Err }

func makeModifyScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyScheduledQueryRequest)

		sq, err := svc.ModifyScheduledQuery(ctx, req.ID, req.payload)
		if err != nil {
			return modifyScheduledQueryResponse{Err: err}, nil
		}

		return modifyScheduledQueryResponse{
			Scheduled: scheduledQueryResponse{
				ScheduledQuery: *sq,
			},
		}, nil
	}
}

type deleteScheduledQueryRequest struct {
	ID uint
}

type deleteScheduledQueryResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteScheduledQueryResponse) error() error { return r.Err }

func makeDeleteScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteScheduledQueryRequest)

		err := svc.DeleteScheduledQuery(ctx, req.ID)
		if err != nil {
			return deleteScheduledQueryResponse{Err: err}, nil
		}

		return deleteScheduledQueryResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) GetScheduledQuery(ctx context.Context, id uint) (*kolide.ScheduledQuery, error) {
	var (
		query *kolide.ScheduledQuery
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetScheduledQuery",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	query, err = mw.Service.GetScheduledQuery(ctx, id)
	return query, err
}

func (mw loggingMiddleware) GetScheduledQueriesInPack(ctx context.Context, id uint, opts kolide.ListOptions) ([]*kolide.ScheduledQuery, error) {
	var (
		queries []*kolide.ScheduledQuery
		err     error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetScheduledQueriesInPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	queries, err = mw.Service.GetScheduledQueriesInPack(ctx, id, opts)
	return queries, err
}

func (mw loggingMiddleware) ScheduleQuery(ctx context.Context, sq *kolide.ScheduledQuery) (*kolide.ScheduledQuery, error) {
	var (
		query *kolide.ScheduledQuery
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ScheduleQuery",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	query, err = mw.Service.ScheduleQuery(ctx, sq)
	return query, err
}

func (mw loggingMiddleware) DeleteScheduledQuery(ctx context.Context, id uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "DeleteScheduledQuery",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.DeleteScheduledQuery(ctx, id)
	return err
}

func (mw loggingMiddleware) ModifyScheduledQuery(ctx context.Context, id uint, p kolide.ScheduledQueryPayload) (*kolide.ScheduledQuery, error) {
	var (
		query *kolide.ScheduledQuery
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ModifyScheduledQuery",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	query, err = mw.Service.ModifyScheduledQuery(ctx, id, p)
	return query, err
}
