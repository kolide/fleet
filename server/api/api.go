package api

import (
	"net/http"

	"github.com/WatchBeam/clock"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/kolide"
)

type service struct {
	ds     kolide.Datastore
	clock  clock.Clock
	logger kitlog.Logger
	jwtKey string
}

func NewService(ds kolide.Datastore, logger kitlog.Logger, jwtKey string, c clock.Clock) (kolide.API, error) {
	return &service{
		ds:     ds,
		logger: logger,
		jwtKey: jwtKey,
		clock:  c,
	}, nil
}

func MakeHandler(apiSvc kolide.API) http.Handler {
	r := mux.NewRouter()

	r.Handle(
		"/api/v2/fleet/hosts",
		makeHttpHandler(
			makeListHostsEndpoint(apiSvc),
			decodeListHostsRequest,
		),
	).Methods("GET").Name("list_hosts")

	return r
}

func makeHttpHandler(e endpoint.Endpoint, decode kithttp.DecodeRequestFunc) http.Handler {
	return kithttp.NewServer(e, decode, encodeResponse)
}
