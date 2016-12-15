package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/datastore/inmem"
	"github.com/kolide/kolide-ose/server/kolide"

	"github.com/stretchr/testify/suite"
)

type EndpointTestSuite struct {
	suite.Suite
	server    *httptest.Server
	userToken string
	ds        kolide.Datastore
}

func (s *EndpointTestSuite) SetupTest() {
	jwtKey := "CHANGEME"
	s.ds, _ = inmem.New(config.TestConfig())
	devOrgInfo := &kolide.AppConfig{
		OrgName:                "Kolide",
		OrgLogoURL:             "http://foo.bar/image.png",
		SMTPPort:               465,
		SMTPAuthenticationType: kolide.AuthTypeUserNamePassword,
		SMTPEnableTLS:          true,
		SMTPVerifySSLCerts:     true,
		SMTPEnableStartTLS:     true,
	}
	s.ds.NewAppConfig(devOrgInfo)
	svc, _ := newTestService(s.ds, nil)
	createTestUsers(s.T(), s.ds)
	logger := kitlog.NewLogfmtLogger(os.Stdout)

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(setRequestsContexts(svc, jwtKey)),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerAfter(kithttp.SetContentType("application/json; charset=utf-8")),
	}

	router := mux.NewRouter()
	ke := MakeKolideServerEndpoints(svc, jwtKey)
	ctxt := context.Background()
	kh := makeKolideKitHandlers(ctxt, ke, opts)
	attachKolideAPIRoutes(router, kh)

	s.server = httptest.NewServer(router)

	userParam := loginRequest{
		Username: "admin1",
		Password: testUsers["admin1"].PlaintextPassword,
	}

	marshalledUser, _ := json.Marshal(&userParam)

	requestBody := &nopCloser{bytes.NewBuffer(marshalledUser)}
	resp, _ := http.Post(s.server.URL+"/api/v1/kolide/login", "application/json", requestBody)

	var jsn = struct {
		User  *kolide.User `json:"user"`
		Token string       `json:"token"`
		Err   string       `json:"error,omitempty"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsn)
	s.userToken = jsn.Token

}

func (s *EndpointTestSuite) TeardownTest() {
	s.server.Close()
}

func TestHttpEndpoints(t *testing.T) {
	suite.Run(t, new(EndpointTestSuite))
}
