package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
)

func makeRequest(server http.Handler, verb, endpoint string, body interface{}, cookie string) *httptest.ResponseRecorder {
	params, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	buff := new(bytes.Buffer)
	buff.Write(params)
	request, _ := http.NewRequest(verb, endpoint, buff)
	if cookie != "" {
		request.Header.Set("Cookie", cookie)
	}
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	return response
}

func createTestServer(ds datastore.Datastore) http.Handler {
	return CreateServer(
		ds,
		kolide.NewMockSMTPConnectionPool(),
		os.Stderr,
		&MockOsqueryResultHandler{},
		&MockOsqueryStatusHandler{},
	)
}

func createTestDatastore() datastore.Datastore {
	ds, err := datastore.New("gorm-sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	return ds
}
