package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/sso"
	"github.com/y0ssar1an/q"
)

////////////////////////////////////////////////////////////////////////////////
// Login
////////////////////////////////////////////////////////////////////////////////

type loginRequest struct {
	Username string // can be username or email
	Password string
}

type loginResponse struct {
	User  *kolide.User `json:"user,omitempty"`
	Token string       `json:"token,omitempty"`
	Err   error        `json:"error,omitempty"`
}

func (r loginResponse) error() error { return r.Err }

func makeLoginEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		user, token, err := svc.Login(ctx, req.Username, req.Password)
		if err != nil {
			return loginResponse{Err: err}, nil
		}
		return loginResponse{user, token, nil}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Logout
////////////////////////////////////////////////////////////////////////////////

type logoutResponse struct {
	Err error `json:"error,omitempty"`
}

func (r logoutResponse) error() error { return r.Err }

func makeLogoutEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := svc.Logout(ctx)
		if err != nil {
			return logoutResponse{Err: err}, nil
		}
		return logoutResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Info About Session
////////////////////////////////////////////////////////////////////////////////

type getInfoAboutSessionRequest struct {
	ID uint
}

type getInfoAboutSessionResponse struct {
	SessionID uint      `json:"session_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Err       error     `json:"error,omitempty"`
}

func (r getInfoAboutSessionResponse) error() error { return r.Err }

func makeGetInfoAboutSessionEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getInfoAboutSessionRequest)
		session, err := svc.GetInfoAboutSession(ctx, req.ID)
		if err != nil {
			return getInfoAboutSessionResponse{Err: err}, nil
		}

		return getInfoAboutSessionResponse{
			SessionID: session.ID,
			UserID:    session.UserID,
			CreatedAt: session.CreatedAt,
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Info About Sessions For User
////////////////////////////////////////////////////////////////////////////////

type getInfoAboutSessionsForUserRequest struct {
	ID uint
}

type getInfoAboutSessionsForUserResponse struct {
	Sessions []getInfoAboutSessionResponse `json:"sessions"`
	Err      error                         `json:"error,omitempty"`
}

func (r getInfoAboutSessionsForUserResponse) error() error { return r.Err }

func makeGetInfoAboutSessionsForUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getInfoAboutSessionsForUserRequest)
		sessions, err := svc.GetInfoAboutSessionsForUser(ctx, req.ID)
		if err != nil {
			return getInfoAboutSessionsForUserResponse{Err: err}, nil
		}
		var resp getInfoAboutSessionsForUserResponse
		for _, session := range sessions {
			resp.Sessions = append(resp.Sessions, getInfoAboutSessionResponse{
				SessionID: session.ID,
				UserID:    session.UserID,
				CreatedAt: session.CreatedAt,
			})
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Session
////////////////////////////////////////////////////////////////////////////////

type deleteSessionRequest struct {
	ID uint
}

type deleteSessionResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteSessionResponse) error() error { return r.Err }

func makeDeleteSessionEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteSessionRequest)
		err := svc.DeleteSession(ctx, req.ID)
		if err != nil {
			return deleteSessionResponse{Err: err}, nil
		}
		return deleteSessionResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Sessions For User
////////////////////////////////////////////////////////////////////////////////

type deleteSessionsForUserRequest struct {
	ID uint
}

type deleteSessionsForUserResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteSessionsForUserResponse) error() error { return r.Err }

func makeDeleteSessionsForUserEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteSessionsForUserRequest)
		err := svc.DeleteSessionsForUser(ctx, req.ID)
		if err != nil {
			return deleteSessionsForUserResponse{Err: err}, nil
		}
		return deleteSessionsForUserResponse{}, nil
	}
}

type initiateSSORequest struct {
	IdentityProviderID uint   `json:"identity_provider_id"`
	RelayURL           string `json:"relay_url"`
	Token              string `json:"token"`
}

type initiateSSOResponse struct {
	URL string `json:"url,omitempty"`
	Err error  `json:"error,omitempty"`
}

func (r initiateSSOResponse) error() error { return r.Err }

func makeInitiateSSOEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		q.Q("endpoint")
		req := request.(initiateSSORequest)
		idProviderURL, err := svc.InitiateSSO(ctx, req.IdentityProviderID, req.RelayURL, req.Token)
		if err != nil {
			return initiateSSOResponse{Err: err}, nil
		}
		return initiateSSOResponse{URL: idProviderURL}, nil
	}
}

func makeLoginSSOEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

type callbackSSOResponse struct {
	URL string `json:"url,omitempty"`
	Err error  `json:"error,omitempty"`
}

func (r callbackSSOResponse) error() error { return r.Err }

// if redirect is present when we encode our response we will
// redirect (302) to this URL
func (r callbackSSOResponse) redirect() string { return r.URL }

func makeCallbackSSOEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		authResponse := request.(sso.AuthInfo)
		// if these two elements are not present they'll be handled in the validation
		// middleware
		userID, _ := authResponse.UserID()
		ssoHandle, _ := authResponse.RelayState()
		redirectURL, err := svc.CallbackSSO(ctx, ssoHandle, userID)
		if err != nil {
			return callbackSSOResponse{Err: err}, nil
		}
		return callbackSSOResponse{URL: redirectURL}, nil
	}
}
