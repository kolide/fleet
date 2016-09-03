package kitserver

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeCreateQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	return req, nil
}