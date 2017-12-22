package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) ListPacks(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Pack, error) {
	return svc.ds.ListPacks(opt)
}

func (svc service) GetPack(ctx context.Context, id uint) (*kolide.Pack, error) {
	return svc.ds.Pack(id)
}

func (svc service) NewPack(ctx context.Context, p kolide.PackPayload) (*kolide.Pack, error) {
	var pack kolide.Pack

	if p.Name != nil {
		pack.Name = *p.Name
	}

	if p.Description != nil {
		pack.Description = *p.Description
	}

	if p.Platform != nil {
		pack.Platform = *p.Platform
	}

	if p.Disabled != nil {
		pack.Disabled = *p.Disabled
	}

	vc, ok := viewer.FromContext(ctx)
	if ok {
		if createdBy := vc.UserID(); createdBy != uint(0) {
			pack.CreatedBy = createdBy
		}
	}

	_, err := svc.ds.NewPack(&pack)
	if err != nil {
		return nil, err
	}

	if p.HostIDs != nil {
		for _, hostID := range *p.HostIDs {
			err = svc.AddHostToPack(ctx, hostID, pack.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	if p.LabelIDs != nil {
		for _, labelID := range *p.LabelIDs {
			err = svc.AddLabelToPack(ctx, labelID, pack.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return &pack, nil
}

func (svc service) ModifyPack(ctx context.Context, id uint, p kolide.PackPayload) (*kolide.Pack, error) {
	pack, err := svc.ds.Pack(id)
	if err != nil {
		return nil, err
	}

	if p.Name != nil {
		pack.Name = *p.Name
	}

	if p.Description != nil {
		pack.Description = *p.Description
	}

	if p.Platform != nil {
		pack.Platform = *p.Platform
	}

	if p.Disabled != nil {
		pack.Disabled = *p.Disabled
	}

	err = svc.ds.SavePack(pack)
	if err != nil {
		return nil, err
	}

	// we must determine what hosts are attached to this pack. then, given
	// our new set of host_ids, we will mutate the database to reflect the
	// desired state.
	if p.HostIDs != nil {

		// first, let's retrieve the total set of hosts
		hosts, err := svc.ListHostsInPack(ctx, pack.ID, kolide.ListOptions{})
		if err != nil {
			return nil, err
		}

		// it will be efficient to create a data structure with constant time
		// lookups to determine whether or not a host is already added
		existingHosts := map[uint]bool{}
		for _, host := range hosts {
			existingHosts[host] = true
		}

		// we will also make a constant time lookup map for the desired set of
		// hosts as well.
		desiredHosts := map[uint]bool{}
		for _, hostID := range *p.HostIDs {
			desiredHosts[hostID] = true
		}

		// if the request declares a host ID but the host is not already
		// associated with the pack, we add it
		for _, hostID := range *p.HostIDs {
			if !existingHosts[hostID] {
				err = svc.AddHostToPack(ctx, hostID, pack.ID)
				if err != nil {
					return nil, err
				}
			}
		}

		// if the request does not declare the ID of a host which currently
		// exists, we delete the existing relationship
		for hostID := range existingHosts {
			if !desiredHosts[hostID] {
				err = svc.RemoveHostFromPack(ctx, hostID, pack.ID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// we must determine what labels are attached to this pack. then, given
	// our new set of label_ids, we will mutate the database to reflect the
	// desired state.
	if p.LabelIDs != nil {

		// first, let's retrieve the total set of labels
		labels, err := svc.ListLabelsForPack(ctx, pack.ID)
		if err != nil {
			return nil, err
		}

		// it will be efficient to create a data structure with constant time
		// lookups to determine whether or not a label is already added
		existingLabels := map[uint]bool{}
		for _, label := range labels {
			existingLabels[label.ID] = true
		}

		// we will also make a constant time lookup map for the desired set of
		// labels as well.
		desiredLabels := map[uint]bool{}
		for _, labelID := range *p.LabelIDs {
			desiredLabels[labelID] = true
		}

		// if the request declares a label ID but the label is not already
		// associated with the pack, we add it
		for _, labelID := range *p.LabelIDs {
			if !existingLabels[labelID] {
				err = svc.AddLabelToPack(ctx, labelID, pack.ID)
				if err != nil {
					return nil, err
				}
			}
		}

		// if the request does not declare the ID of a label which currently
		// exists, we delete the existing relationship
		for labelID := range existingLabels {
			if !desiredLabels[labelID] {
				err = svc.RemoveLabelFromPack(ctx, labelID, pack.ID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return pack, err
}

func (svc service) DeletePack(ctx context.Context, id uint) error {
	return svc.ds.DeletePack(id)
}

func (svc service) AddLabelToPack(ctx context.Context, lid, pid uint) error {
	return svc.ds.AddLabelToPack(lid, pid)
}

func (svc service) ListLabelsForPack(ctx context.Context, pid uint) ([]*kolide.Label, error) {
	return svc.ds.ListLabelsForPack(pid)
}

func (svc service) RemoveLabelFromPack(ctx context.Context, lid, pid uint) error {
	return svc.ds.RemoveLabelFromPack(lid, pid)
}

func (svc service) AddHostToPack(ctx context.Context, hid, pid uint) error {
	return svc.ds.AddHostToPack(hid, pid)
}

func (svc service) RemoveHostFromPack(ctx context.Context, hid, pid uint) error {
	return svc.ds.RemoveHostFromPack(hid, pid)
}

func (svc service) ListHostsInPack(ctx context.Context, pid uint, opt kolide.ListOptions) ([]uint, error) {
	return svc.ds.ListHostsInPack(pid, opt)
}

func (svc service) ListExplicitHostsInPack(ctx context.Context, pid uint, opt kolide.ListOptions) ([]uint, error) {
	return svc.ds.ListExplicitHostsInPack(pid, opt)
}

func (svc service) ListPacksForHost(ctx context.Context, hid uint) ([]*kolide.Pack, error) {
	packs := []*kolide.Pack{}

	// we will need to give some subset of packs to this host based on the
	// labels which this host is known to belong to
	allPacks, err := svc.ds.ListPacks(kolide.ListOptions{})
	if err != nil {
		return nil, err
	}

	// pull the labels that this host belongs to
	labels, err := svc.ds.ListLabelsForHost(hid)
	if err != nil {
		return nil, err
	}

	// in order to use o(1) array indexing in an o(n) loop vs a o(n^2) double
	// for loop iteration, we must create the array which may be indexed below
	labelIDs := map[uint]bool{}
	for _, label := range labels {
		labelIDs[label.ID] = true
	}

	for _, pack := range allPacks {
		// don't include packs which have been disabled
		if pack.Disabled {
			continue
		}

		// for each pack, we must know what labels have been assigned to that
		// pack
		labelsForPack, err := svc.ds.ListLabelsForPack(pack.ID)
		if err != nil {
			return nil, err
		}

		// o(n) iteration to determine whether or not a pack is enabled
		// in this case, n is len(labelsForPack)
		for _, label := range labelsForPack {
			if labelIDs[label.ID] {
				packs = append(packs, pack)
				break
			}
		}

		// for each pack, we must know what host have been assigned to that pack
		hostsForPack, err := svc.ds.ListExplicitHostsInPack(pack.ID, kolide.ListOptions{})
		if err != nil {
			return nil, err
		}

		// o(n) iteration to determine whether or not a pack is enabled
		// in this case, n is len(hostsForPack)
		for _, host := range hostsForPack {
			if host == hid {
				packs = append(packs, pack)
				break
			}
		}
	}

	return packs, nil
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeCreatePackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeModifyPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeDeletePackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req deletePackRequest
	req.ID = id
	return req, nil
}

func decodeGetPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getPackRequest
	req.ID = id
	return req, nil
}

func decodeListPacksRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listPacksRequest{ListOptions: opt}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type packResponse struct {
	kolide.Pack
	QueryCount uint `json:"query_count"`

	// All current hosts in the pack. Hosts which are selected explicty and
	// hosts which are part of a label.
	TotalHostsCount uint `json:"total_hosts_count"`

	// IDs of hosts which were explicitly selected.
	HostIDs  []uint `json:"host_ids"`
	LabelIDs []uint `json:"label_ids"`
}

func packResponseForPack(ctx context.Context, svc kolide.Service, pack kolide.Pack) (*packResponse, error) {
	opts := kolide.ListOptions{}
	queries, err := svc.GetScheduledQueriesInPack(ctx, pack.ID, opts)
	if err != nil {
		return nil, err
	}

	hosts, err := svc.ListExplicitHostsInPack(ctx, pack.ID, opts)
	if err != nil {
		return nil, err
	}

	labels, err := svc.ListLabelsForPack(ctx, pack.ID)
	labelIDs := make([]uint, len(labels))
	for i, label := range labels {
		labelIDs[i] = label.ID
	}
	if err != nil {
		return nil, err
	}

	hostMetrics, err := svc.CountHostsInTargets(ctx, hosts, labelIDs)
	if err != nil {
		return nil, err
	}

	return &packResponse{
		Pack:            pack,
		QueryCount:      uint(len(queries)),
		TotalHostsCount: hostMetrics.TotalHosts,
		HostIDs:         hosts,
		LabelIDs:        labelIDs,
	}, nil
}

type getPackRequest struct {
	ID uint
}

type getPackResponse struct {
	Pack packResponse `json:"pack,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r getPackResponse) error() error { return r.Err }

func makeGetPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPackRequest)

		pack, err := svc.GetPack(ctx, req.ID)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		resp, err := packResponseForPack(ctx, svc, *pack)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		return getPackResponse{
			Pack: *resp,
		}, nil
	}
}

type listPacksRequest struct {
	ListOptions kolide.ListOptions
}

type listPacksResponse struct {
	Packs []packResponse `json:"packs"`
	Err   error          `json:"error,omitempty"`
}

func (r listPacksResponse) error() error { return r.Err }

func makeListPacksEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listPacksRequest)
		packs, err := svc.ListPacks(ctx, req.ListOptions)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		resp := listPacksResponse{Packs: make([]packResponse, len(packs))}
		for i, pack := range packs {
			packResp, err := packResponseForPack(ctx, svc, *pack)
			if err != nil {
				return getPackResponse{Err: err}, nil
			}
			resp.Packs[i] = *packResp
		}
		return resp, nil
	}
}

type createPackRequest struct {
	payload kolide.PackPayload
}

type createPackResponse struct {
	Pack packResponse `json:"pack,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r createPackResponse) error() error { return r.Err }

func makeCreatePackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createPackRequest)
		pack, err := svc.NewPack(ctx, req.payload)
		if err != nil {
			return createPackResponse{Err: err}, nil
		}

		resp, err := packResponseForPack(ctx, svc, *pack)
		if err != nil {
			return createPackResponse{Err: err}, nil
		}

		return createPackResponse{
			Pack: *resp,
		}, nil
	}
}

type modifyPackRequest struct {
	ID      uint
	payload kolide.PackPayload
}

type modifyPackResponse struct {
	Pack packResponse `json:"pack,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r modifyPackResponse) error() error { return r.Err }

func makeModifyPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyPackRequest)
		pack, err := svc.ModifyPack(ctx, req.ID, req.payload)
		if err != nil {
			return modifyPackResponse{Err: err}, nil
		}

		resp, err := packResponseForPack(ctx, svc, *pack)
		if err != nil {
			return modifyPackResponse{Err: err}, nil
		}

		return modifyPackResponse{
			Pack: *resp,
		}, nil
	}
}

type deletePackRequest struct {
	ID uint
}

type deletePackResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deletePackResponse) error() error { return r.Err }

func makeDeletePackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePackRequest)
		err := svc.DeletePack(ctx, req.ID)
		if err != nil {
			return deletePackResponse{Err: err}, nil
		}
		return deletePackResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) ListPacks(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Pack, error) {
	var (
		packs []*kolide.Pack
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListPacks",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	packs, err = mw.Service.ListPacks(ctx, opt)
	return packs, err
}

func (mw loggingMiddleware) GetPack(ctx context.Context, id uint) (*kolide.Pack, error) {
	var (
		pack *kolide.Pack
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	pack, err = mw.Service.GetPack(ctx, id)
	return pack, err
}

func (mw loggingMiddleware) NewPack(ctx context.Context, p kolide.PackPayload) (*kolide.Pack, error) {
	var (
		pack *kolide.Pack
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "NewPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	pack, err = mw.Service.NewPack(ctx, p)
	return pack, err
}

func (mw loggingMiddleware) ModifyPack(ctx context.Context, id uint, p kolide.PackPayload) (*kolide.Pack, error) {
	var (
		pack *kolide.Pack
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ModifyPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	pack, err = mw.Service.ModifyPack(ctx, id, p)
	return pack, err
}

func (mw loggingMiddleware) DeletePack(ctx context.Context, id uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "DeletePack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.DeletePack(ctx, id)
	return err
}

func (mw loggingMiddleware) AddLabelToPack(ctx context.Context, lid uint, pid uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "AddLabelToPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.AddLabelToPack(ctx, lid, pid)
	return err
}

func (mw loggingMiddleware) RemoveLabelFromPack(ctx context.Context, lid uint, pid uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "RemoveLabelFromPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.RemoveLabelFromPack(ctx, lid, pid)
	return err
}

func (mw loggingMiddleware) ListLabelsForPack(ctx context.Context, pid uint) ([]*kolide.Label, error) {
	var (
		labels []*kolide.Label
		err    error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListLabelsForPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	labels, err = mw.Service.ListLabelsForPack(ctx, pid)
	return labels, err
}

func (mw loggingMiddleware) AddHostToPack(ctx context.Context, hid uint, pid uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "AddHostToPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.AddHostToPack(ctx, hid, pid)
	return err
}

func (mw loggingMiddleware) RemoveHostFromPack(ctx context.Context, hid uint, pid uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "RemoveHostFromPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.RemoveHostFromPack(ctx, hid, pid)
	return err
}

func (mw loggingMiddleware) ListPacksForHost(ctx context.Context, hid uint) ([]*kolide.Pack, error) {
	var (
		packs []*kolide.Pack
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListPacksForHost",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	packs, err = mw.Service.ListPacksForHost(ctx, hid)
	return packs, err
}

func (mw loggingMiddleware) ListHostsInPack(ctx context.Context, pid uint, opt kolide.ListOptions) ([]uint, error) {
	var (
		hosts []uint
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListHostsInPack",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	hosts, err = mw.Service.ListHostsInPack(ctx, pid, opt)
	return hosts, err
}
