package service

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeNewIdentityProviderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req newIdentityProviderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeModifyIdentityProviderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyIdentityProviderRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.id = id
	return req, nil
}

func decodeGetIdentityProviderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getIdentityProviderRequest{id}, nil
}

func decodeDeleteIdentityProviderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteIdentityProviderRequest{id}, nil
}
