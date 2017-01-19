package service

import (
	"net/http"
	"strconv"

	"golang.org/x/net/context"
)

func decodeCertificateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var insecure bool
	v := r.URL.Query().Get("insecure")
	if s, err := strconv.ParseBool(v); err == nil {
		insecure = s
	}
	return certificateRequest{Insecure: insecure}, nil
}
