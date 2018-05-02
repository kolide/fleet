package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

type SetupRequest struct {
	Admin           *kolide.UserPayload `json:"admin"`
	OrgInfo         *kolide.OrgInfo     `json:"org_info"`
	KolideServerURL *string             `json:"kolide_server_url,omitempty"`
	EnrollSecret    *string             `json:"osquery_enroll_secret,omitempty"`
}

type SetupResponse struct {
	Admin           *kolide.User    `json:"admin,omitempty"`
	OrgInfo         *kolide.OrgInfo `json:"org_info,omitempty"`
	KolideServerURL *string         `json:"kolide_server_url"`
	EnrollSecret    *string         `json:"osquery_enroll_secret"`
	Token           *string         `json:"token,omitempty"`
	Err             error           `json:"error,omitempty"`
}

func (r SetupResponse) error() error { return r.Err }

func makeSetupEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			admin         *kolide.User
			config        *kolide.AppConfig
			configPayload kolide.AppConfigPayload
			err           error
		)
		req := request.(SetupRequest)
		if req.OrgInfo != nil {
			configPayload.OrgInfo = req.OrgInfo
		}
		configPayload.ServerSettings = &kolide.ServerSettings{}
		if req.KolideServerURL != nil {
			configPayload.ServerSettings.KolideServerURL = req.KolideServerURL
		}
		if req.EnrollSecret != nil {
			configPayload.ServerSettings.EnrollSecret = req.EnrollSecret
		}
		config, err = svc.NewAppConfig(ctx, configPayload)
		if err != nil {
			return SetupResponse{Err: err}, nil
		}
		// creating the user should be the last action. If there's a user
		// present and other errors occur, the setup endpoint closes.
		if req.Admin != nil {
			admin, err = svc.NewAdminCreatedUser(ctx, *req.Admin)
			if err != nil {
				return SetupResponse{Err: err}, nil
			}
		}

		// If everything works to this point, log the user in and return token.  If
		// the login fails for some reason, ignore the error and don't return
		// a token, forcing the user to log in manually
		token := new(string)
		_, *token, err = svc.Login(ctx, *req.Admin.Username, *req.Admin.Password)
		if err != nil {
			token = nil
		}
		return SetupResponse{
			Admin: admin,
			OrgInfo: &kolide.OrgInfo{
				OrgName:    &config.OrgName,
				OrgLogoURL: &config.OrgLogoURL,
			},
			KolideServerURL: &config.KolideServerURL,
			EnrollSecret:    &config.EnrollSecret,
			Token:           token,
		}, nil
	}
}
