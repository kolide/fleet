package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/adams-sarah/test2doc/test"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/kolide/kolide-ose/server/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIDoc(t *testing.T) {
	ds, err := datastore.New("inmem", "")
	assert.Nil(t, err)
	createTestUsers(t, ds)
	svc, err := newTestService(ds)
	assert.Nil(t, err)
	handler := MakeHandler(context.Background(), svc, "CHANGEME", kitlog.NewNopLogger())
	test.RegisterURLVarExtractor(mux.Vars)
	server, err := test.NewServer(handler)
	assert.Nil(t, err)
	defer server.Finish()
	admin1JWT := loggedInSessionToken(t, handler)

	var docTests = []struct {
		path         string
		method       string
		body         io.Reader
		sessionToken string
	}{
		{
			path:   "/api/v1/kolide/login",
			method: "POST",
			body: newBody([]byte(`{
        		"username": "admin1",
        		"password": "foobar"
			}`)),
		},
		{
			path:   "/api/v1/kolide/login",
			method: "POST",
			body: newBody([]byte(`{
        		"username": "nobody",
        		"password": "foobar"
			}`)),
		},
		{
			path:         "/api/v1/kolide/logout",
			method:       "POST",
			sessionToken: admin1JWT,
		},
	}

	client := http.DefaultClient
	for _, tt := range docTests {
		req := httptest.NewRequest(tt.method, tt.path, tt.body)
		if tt.sessionToken != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.sessionToken))
		}
		req.RequestURI = ""
		req.URL, err = url.Parse(fmt.Sprintf("%s%s", server.URL, req.URL.Path))
		require.Nil(t, err)
		_, err = client.Do(req)
		require.Nil(t, err)
	}
}

func newBody(data []byte) *bytes.Buffer {
	var body bytes.Buffer
	body.Write(data)
	return &body
}

func loggedInSessionToken(t *testing.T, handler http.Handler) string {
	var userToken struct {
		Token string
	}
	request := httptest.NewRequest("POST", "/api/v1/kolide/login", newBody([]byte(`{
        		"username": "admin1",
        		"password": "foobar"
			}`)))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	err := json.NewDecoder(response.Body).Decode(&userToken)
	require.Nil(t, err)
	return userToken.Token
}
