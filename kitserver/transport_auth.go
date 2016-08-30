package kitserver

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return loginResponse{Err: err}, nil
	}
	defer r.Body.Close()
	return req, nil
}
