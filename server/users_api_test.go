package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
	"github.com/kolide/kolide-ose/server"
)

func TestCreateUserAPI(t *testing.T) {
	ds, _ := datastore.New("mock", "")
	// svc, _ := kolide.NewService(ds)
	srv := newServer(t, ds)
	defer srv.Close()

	var createUserTests = []struct {
		Username string
		Password string
	}{
		{Username: "admin1",
			Password: "foobar",
		},
	}

	for _, tt := range createUserTests {
		resp := do(t, srv, "/api/v1/kolide/login", "POST", server.CreateUserRequestBody{
			Username: tt.Username,
			Password: tt.Password,
		})

		// userCookie := resp.Header.Get("Set-Cookie")
		// if userCookie == "" {
		// 	t.Fatal("login handle returned empty cookie")
		// }
		if resp.StatusCode == http.StatusOK {
			t.Errorf("expected OK, got %q", resp.Status)
		}
	}

}

func newServer(t *testing.T, ds kolide.Datastore) *httptest.Server {
	handler := server.CreateServer(
		ds,
		kolide.NewMockSMTPConnectionPool(),
		ioutil.Discard,
		&server.MockOsqueryResultHandler{},
		&server.MockOsqueryStatusHandler{},
	)
	return httptest.NewServer(handler)
}

func do(t *testing.T, srv *httptest.Server, endpoint, method string, jsn interface{}) *http.Response {
	client := http.DefaultClient
	theURL := srv.URL + endpoint
	data, err := json.Marshal(jsn)
	if err != nil {
		t.Fatal(err)
	}
	body := ioutil.NopCloser(bytes.NewBuffer(data))
	resp, err := client.Post(theURL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
