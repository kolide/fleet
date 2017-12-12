package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) ListLabels(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Label, error) {
	return svc.ds.ListLabels(opt)
}

func (svc service) GetLabel(ctx context.Context, id uint) (*kolide.Label, error) {
	return svc.ds.Label(id)
}

func (svc service) NewLabel(ctx context.Context, p kolide.LabelPayload) (*kolide.Label, error) {
	label := &kolide.Label{}

	if p.Name == nil {
		return nil, newInvalidArgumentError("name", "missing required argument")
	}
	label.Name = *p.Name

	if p.Query == nil {
		return nil, newInvalidArgumentError("query", "missing required argument")
	}
	label.Query = *p.Query

	if p.Platform != nil {
		label.Platform = *p.Platform
	}

	if p.Description != nil {
		label.Description = *p.Description
	}

	label, err := svc.ds.NewLabel(label)
	if err != nil {
		return nil, err
	}
	return label, nil
}

func (svc service) DeleteLabel(ctx context.Context, id uint) error {
	return svc.ds.DeleteLabel(id)
}

func (svc service) HostIDsForLabel(lid uint) ([]uint, error) {
	hosts, err := svc.ds.ListHostsInLabel(lid)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, h := range hosts {
		ids = append(ids, h.ID)
	}
	return ids, nil
}

func (svc service) ModifyLabel(ctx context.Context, id uint, payload kolide.ModifyLabelPayload) (*kolide.Label, error) {
	label, err := svc.ds.Label(id)
	if err != nil {
		return nil, err
	}
	if payload.Name != nil {
		label.Name = *payload.Name
	}
	if payload.Description != nil {
		label.Description = *payload.Description
	}
	return svc.ds.SaveLabel(label)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeCreateLabelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeDeleteLabelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req deleteLabelRequest
	req.ID = id
	return req, nil
}

func decodeGetLabelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getLabelRequest
	req.ID = id
	return req, nil
}

func decodeListLabelsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listLabelsRequest{ListOptions: opt}, nil
}

func decodeModifyLabelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var resp modifyLabelRequest
	err = json.NewDecoder(r.Body).Decode(&resp.payload)
	if err != nil {
		return nil, err
	}
	resp.ID = id
	return resp, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type getLabelRequest struct {
	ID uint
}

type labelResponse struct {
	kolide.Label
	DisplayText     string `json:"display_text"`
	Count           uint   `json:"count"`
	Online          uint   `json:"online"`
	Offline         uint   `json:"offline"`
	MissingInAction uint   `json:"missing_in_action"`
	HostIDs         []uint `json:"host_ids"`
}

type getLabelResponse struct {
	Label labelResponse `json:"label"`
	Err   error         `json:"error,omitempty"`
}

func (r getLabelResponse) error() error { return r.Err }

func labelResponseForLabel(ctx context.Context, svc kolide.Service, label *kolide.Label) (*labelResponse, error) {
	metrics, err := svc.CountHostsInTargets(ctx, nil, []uint{label.ID})
	if err != nil {
		return nil, err
	}
	hosts, err := svc.HostIDsForLabel(label.ID)
	if err != nil {
		return nil, err
	}
	return &labelResponse{
		*label,
		label.Name,
		metrics.TotalHosts,
		metrics.OnlineHosts,
		metrics.OfflineHosts,
		metrics.MissingInActionHosts,
		hosts,
	}, nil
}

func makeGetLabelEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getLabelRequest)
		label, err := svc.GetLabel(ctx, req.ID)
		if err != nil {
			return getLabelResponse{Err: err}, nil
		}
		resp, err := labelResponseForLabel(ctx, svc, label)
		if err != nil {
			return getLabelResponse{Err: err}, nil
		}
		return getLabelResponse{Label: *resp}, nil
	}
}

type listLabelsRequest struct {
	ListOptions kolide.ListOptions
}

type listLabelsResponse struct {
	Labels []labelResponse `json:"labels"`
	Err    error           `json:"error,omitempty"`
}

func (r listLabelsResponse) error() error { return r.Err }

func makeListLabelsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listLabelsRequest)
		labels, err := svc.ListLabels(ctx, req.ListOptions)
		if err != nil {
			return listLabelsResponse{Err: err}, nil
		}

		resp := listLabelsResponse{}
		for _, label := range labels {
			labelResp, err := labelResponseForLabel(ctx, svc, label)
			if err != nil {
				return listLabelsResponse{Err: err}, nil
			}
			resp.Labels = append(resp.Labels, *labelResp)
		}
		return resp, nil
	}
}

type createLabelRequest struct {
	payload kolide.LabelPayload
}

type createLabelResponse struct {
	Label labelResponse `json:"label"`
	Err   error         `json:"error,omitempty"`
}

func (r createLabelResponse) error() error { return r.Err }

func makeCreateLabelEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createLabelRequest)

		label, err := svc.NewLabel(ctx, req.payload)
		if err != nil {
			return createLabelResponse{Err: err}, nil
		}

		labelResp, err := labelResponseForLabel(ctx, svc, label)
		if err != nil {
			return createLabelResponse{Err: err}, nil
		}

		return createLabelResponse{Label: *labelResp}, nil
	}
}

type deleteLabelRequest struct {
	ID uint
}

type deleteLabelResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteLabelResponse) error() error { return r.Err }

func makeDeleteLabelEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteLabelRequest)
		err := svc.DeleteLabel(ctx, req.ID)
		if err != nil {
			return deleteLabelResponse{Err: err}, nil
		}
		return deleteLabelResponse{}, nil
	}
}

type modifyLabelRequest struct {
	ID      uint
	payload kolide.ModifyLabelPayload
}

type modifyLabelResponse struct {
	Label labelResponse `json:"label"`
	Err   error         `json:"error,omitempty"`
}

func (r modifyLabelResponse) error() error { return r.Err }

func makeModifyLabelEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyLabelRequest)
		label, err := svc.ModifyLabel(ctx, req.ID, req.payload)
		if err != nil {
			return modifyLabelResponse{Err: err}, nil
		}

		labelResp, err := labelResponseForLabel(ctx, svc, label)
		if err != nil {
			return modifyLabelResponse{Err: err}, nil
		}

		return modifyLabelResponse{Label: *labelResp}, err
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) ModifyLabel(ctx context.Context, id uint, p kolide.ModifyLabelPayload) (*kolide.Label, error) {
	var (
		lic *kolide.Label
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyLabel", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	lic, err = mw.Service.ModifyLabel(ctx, id, p)
	return lic, err

}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) ModifyLabel(ctx context.Context, id uint, p kolide.ModifyLabelPayload) (*kolide.Label, error) {
	var (
		label *kolide.Label
		err   error
	)

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ModifyLabel",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	label, err = mw.Service.ModifyLabel(ctx, id, p)
	return label, err
}

func (mw loggingMiddleware) ListLabels(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Label, error) {
	var (
		labels []*kolide.Label
		err    error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListLabels",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	labels, err = mw.Service.ListLabels(ctx, opt)
	return labels, err
}

func (mw loggingMiddleware) GetLabel(ctx context.Context, id uint) (*kolide.Label, error) {
	var (
		label *kolide.Label
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetLabel",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	label, err = mw.Service.GetLabel(ctx, id)
	return label, err
}

func (mw loggingMiddleware) NewLabel(ctx context.Context, p kolide.LabelPayload) (*kolide.Label, error) {
	var (
		label *kolide.Label
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "NewLabel",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	label, err = mw.Service.NewLabel(ctx, p)
	return label, err
}

func (mw loggingMiddleware) DeleteLabel(ctx context.Context, id uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "DeleteLabel",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.DeleteLabel(ctx, id)
	return err
}
