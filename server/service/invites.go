package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) InviteNewUser(ctx context.Context, payload kolide.InvitePayload) (*kolide.Invite, error) {
	// verify that the user with the given email does not already exist
	_, err := svc.ds.UserByEmail(*payload.Email)
	if err == nil {
		return nil, newInvalidArgumentError("email", "a user with this account already exists")
	}

	if _, ok := err.(kolide.NotFoundError); !ok {
		return nil, err
	}

	// find the user who created the invite
	inviter, err := svc.User(ctx, *payload.InvitedBy)
	if err != nil {
		return nil, err
	}

	random, err := kolide.RandomText(svc.config.App.TokenKeySize)
	if err != nil {
		return nil, err
	}
	token := base64.URLEncoding.EncodeToString([]byte(random))

	invite := &kolide.Invite{
		Email:     *payload.Email,
		Admin:     *payload.Admin,
		InvitedBy: inviter.ID,
		Token:     token,
	}
	if payload.Position != nil {
		invite.Position = *payload.Position
	}
	if payload.Name != nil {
		invite.Name = *payload.Name
	}
	if payload.SSOEnabled != nil {
		invite.SSOEnabled = *payload.SSOEnabled
	}

	invite, err = svc.ds.NewInvite(invite)
	if err != nil {
		return nil, err
	}

	config, err := svc.AppConfig(ctx)
	if err != nil {
		return nil, err
	}

	invitedBy := inviter.Name
	if invitedBy == "" {
		invitedBy = inviter.Username
	}
	inviteEmail := kolide.Email{
		Subject: "You're Invited to Kolide",
		To:      []string{invite.Email},
		Config:  config,
		Mailer: &kolide.InviteMailer{
			Invite:            invite,
			KolideServerURL:   template.URL(config.KolideServerURL),
			OrgName:           config.OrgName,
			InvitedByUsername: invitedBy,
		},
	}

	err = svc.mailService.SendEmail(inviteEmail)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

func (svc service) ListInvites(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Invite, error) {
	return svc.ds.ListInvites(opt)
}

func (svc service) VerifyInvite(ctx context.Context, token string) (*kolide.Invite, error) {
	invite, err := svc.ds.InviteByToken(token)
	if err != nil {
		return nil, err
	}

	if invite.Token != token {
		return nil, newInvalidArgumentError("invite_token", "Invite Token does not match Email Address.")
	}

	expiresAt := invite.CreatedAt.Add(svc.config.App.InviteTokenValidityPeriod)
	if svc.clock.Now().After(expiresAt) {
		return nil, newInvalidArgumentError("invite_token", "Invite token has expired.")
	}

	return invite, nil

}

func (svc service) DeleteInvite(ctx context.Context, id uint) error {
	return svc.ds.DeleteInvite(id)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeCreateInviteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	if req.payload.Email != nil {
		*req.payload.Email = strings.ToLower(*req.payload.Email)
	}

	return req, nil
}

func decodeDeleteInviteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteInviteRequest{ID: id}, nil
}

func decodeVerifyInviteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		return 0, errBadRoute
	}
	return verifyInviteRequest{Token: token}, nil
}

func decodeListInvitesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listInvitesRequest{ListOptions: opt}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type createInviteRequest struct {
	payload kolide.InvitePayload
}

type createInviteResponse struct {
	Invite *kolide.Invite `json:"invite,omitempty"`
	Err    error          `json:"error,omitempty"`
}

func (r createInviteResponse) error() error { return r.Err }

func makeCreateInviteEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createInviteRequest)
		invite, err := svc.InviteNewUser(ctx, req.payload)
		if err != nil {
			return createInviteResponse{Err: err}, nil
		}
		return createInviteResponse{invite, nil}, nil
	}
}

type listInvitesRequest struct {
	ListOptions kolide.ListOptions
}

type listInvitesResponse struct {
	Invites []kolide.Invite `json:"invites"`
	Err     error           `json:"error,omitempty"`
}

func (r listInvitesResponse) error() error { return r.Err }

func makeListInvitesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listInvitesRequest)
		invites, err := svc.ListInvites(ctx, req.ListOptions)
		if err != nil {
			return listInvitesResponse{Err: err}, nil
		}

		resp := listInvitesResponse{Invites: []kolide.Invite{}}
		for _, invite := range invites {
			resp.Invites = append(resp.Invites, *invite)
		}
		return resp, nil
	}
}

type deleteInviteRequest struct {
	ID uint
}

type deleteInviteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteInviteResponse) error() error { return r.Err }

func makeDeleteInviteEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteInviteRequest)
		err := svc.DeleteInvite(ctx, req.ID)
		if err != nil {
			return deleteInviteResponse{Err: err}, nil
		}
		return deleteInviteResponse{}, nil
	}
}

type verifyInviteRequest struct {
	Token string
}

type verifyInviteResponse struct {
	Invite *kolide.Invite `json:"invite"`
	Err    error          `json:"error,omitempty"`
}

func (r verifyInviteResponse) error() error { return r.Err }

func makeVerifyInviteEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyInviteRequest)
		invite, err := svc.VerifyInvite(ctx, req.Token)
		if err != nil {
			return verifyInviteResponse{Err: err}, nil
		}
		return verifyInviteResponse{Invite: invite}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) InviteNewUser(ctx context.Context, payload kolide.InvitePayload) (*kolide.Invite, error) {
	var (
		invite *kolide.Invite
		err    error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "InviteNewUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	invite, err = mw.Service.InviteNewUser(ctx, payload)
	return invite, err
}

func (mw metricsMiddleware) DeleteInvite(ctx context.Context, id uint) error {
	var (
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteInvite", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.DeleteInvite(ctx, id)
	return err
}

func (mw metricsMiddleware) ListInvites(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Invite, error) {
	var (
		invites []*kolide.Invite
		err     error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "Invites", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	invites, err = mw.Service.ListInvites(ctx, opt)
	return invites, err
}

func (mw metricsMiddleware) VerifyInvite(ctx context.Context, token string) (*kolide.Invite, error) {
	var (
		err    error
		invite *kolide.Invite
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "VerifyInvite", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	invite, err = mw.Service.VerifyInvite(ctx, token)
	return invite, err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) InviteNewUser(ctx context.Context, payload kolide.InvitePayload) (*kolide.Invite, error) {
	var (
		invite *kolide.Invite
		err    error
	)

	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "InviteNewUser",
			"created_by", vc.Username(),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	invite, err = mw.Service.InviteNewUser(ctx, payload)
	return invite, err
}

func (mw loggingMiddleware) DeleteInvite(ctx context.Context, id uint) error {
	var (
		err error
	)
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return errNoContext
	}
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "DeleteInvite",
			"deleted_by", vc.Username(),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Service.DeleteInvite(ctx, id)
	return err
}

func (mw loggingMiddleware) ListInvites(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Invite, error) {
	var (
		invites []*kolide.Invite
		err     error
	)
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Invites",
			"called_by", vc.Username(),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	invites, err = mw.Service.ListInvites(ctx, opt)
	return invites, err
}

func (mw loggingMiddleware) VerifyInvite(ctx context.Context, token string) (*kolide.Invite, error) {
	var (
		err    error
		invite *kolide.Invite
	)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "VerifyInvite",
			"token", token,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	invite, err = mw.Service.VerifyInvite(ctx, token)
	return invite, err
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func (mw validationMiddleware) InviteNewUser(ctx context.Context, payload kolide.InvitePayload) (*kolide.Invite, error) {
	invalid := &invalidArgumentError{}
	if payload.Email == nil {
		invalid.Append("email", "missing required argument")
	}
	if payload.InvitedBy == nil {
		invalid.Append("invited_by", "missing required argument")
	}
	if payload.Admin == nil {
		invalid.Append("admin", "missing required argument")
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.InviteNewUser(ctx, payload)
}
