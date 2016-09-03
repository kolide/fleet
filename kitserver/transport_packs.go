package kitserver

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeCreatePackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeModifyPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r)
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
