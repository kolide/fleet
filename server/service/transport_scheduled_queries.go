package service

import (
	"net/http"

	"golang.org/x/net/context"
)

func decodeScheduleQueriesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req scheduleQueriesRequest

	// TODO: parse request body into req
	// if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
	//	return nil, err
	// }

	return req, nil
}

func decodeModifyScheduledQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyScheduledQueryRequest

	// TODO: parse request body into req
	// if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
	//	return nil, err
	// }

	req.ID = id
	return req, nil
}

func decodeDeleteScheduledQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req deleteScheduledQueryRequest
	req.ID = id
	return req, nil
}

func decodeGetScheduledQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getScheduledQueryRequest
	req.ID = id
	return req, nil
}

func decodeGetScheduledQueriesInPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getScheduledQueriesInPackRequest
	req.ID = id
	return req, nil
}
