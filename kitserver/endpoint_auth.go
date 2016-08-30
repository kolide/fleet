package kitserver

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/kolide"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// for now they are one and the same
type loginResponse getUserResponse

// this endpoint is for discussion on the PR.
// IMO we should just have a regular `/login` http.Handler for login which accepts the service
// as an argument
// go-kit is great for stateless request/response and the SessionManager
// interacts with the reader/writer to manage state. Wiring it as a go-kit endpoint is awkward since
// the ResponseWriter and the Request are never in the same place.
func makeLoginEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		user, err := svc.Authenticate(ctx, req.Username, req.Password)
		if err != nil {
			return loginResponse{Err: err}, nil
		}

		return loginResponse{
			ID:                 user.ID,
			Username:           user.Username,
			Email:              user.Email,
			Admin:              user.Admin,
			Enabled:            user.Enabled,
			NeedsPasswordReset: user.NeedsPasswordReset,
		}, nil
	}
}
