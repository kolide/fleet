package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Get Label
////////////////////////////////////////////////////////////////////////////////

type getLabelRequest struct {
	ID uint
}

type labelResponse struct {
	kolide.Label
	Count uint `json:"count"`
}

type getLabelResponse struct {
	Label labelResponse `json:"label"`
	Err   error         `json:"error,omitempty"`
}

func (r getLabelResponse) error() error { return r.Err }

func makeGetLabelEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getLabelRequest)
		label, err := svc.GetLabel(ctx, req.ID)
		if err != nil {
			return getLabelResponse{Err: err}, nil
		}
		count, err := svc.CountHostsInTargets(ctx, nil, []uint{label.ID})
		if err != nil {
			return getLabelResponse{Err: err}, nil
		}
		return getLabelResponse{
			labelResponse{
				*label,
				count,
			},
			nil,
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// List Labels
////////////////////////////////////////////////////////////////////////////////

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

		resp := listLabelsResponse{Labels: []labelResponse{}}
		for _, label := range labels {
			count, err := svc.CountHostsInTargets(ctx, nil, []uint{label.ID})
			if err != nil {
				return listLabelsResponse{Err: err}, nil
			}
			resp.Labels = append(resp.Labels,
				labelResponse{
					*label,
					count,
				},
			)
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Create Label
////////////////////////////////////////////////////////////////////////////////

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
		count, err := svc.CountHostsInTargets(ctx, nil, []uint{label.ID})
		if err != nil {
			return createLabelResponse{Err: err}, nil
		}
		return createLabelResponse{
			labelResponse{
				*label,
				count,
			},
			nil,
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify Label
////////////////////////////////////////////////////////////////////////////////

type modifyLabelRequest struct {
	ID      uint
	payload kolide.LabelPayload
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
		count, err := svc.CountHostsInTargets(ctx, nil, []uint{label.ID})
		if err != nil {
			return modifyLabelResponse{Err: err}, nil
		}

		return modifyLabelResponse{
			labelResponse{
				*label,
				count,
			},
			nil,
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Label
////////////////////////////////////////////////////////////////////////////////

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
