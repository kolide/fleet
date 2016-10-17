package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/adams-sarah/test2doc/test"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/kolide/kolide-ose/server/datastore"
	"github.com/stretchr/testify/assert"
)

func TestAPIDoc(t *testing.T) {
	ds, err := datastore.New("inmem", "")
	assert.Nil(t, err)
	createTestUsers(t, ds)
	svc, err := newTestService(ds)
	assert.Nil(t, err)

	handler := MakeHandler(context.Background(), svc, "CHANEME", kitlog.NewNopLogger())
	test.RegisterURLVarExtractor(mux.Vars)
	server, err := test.NewServer(handler)
	assert.Nil(t, err)
	defer server.Finish()

	var body bytes.Buffer
	body.Write([]byte(`{
        "username": "admin1",
        "password": "foobar"
    }`))
	req := httptest.NewRequest("POST", "/api/v1/kolide/login", &body)
	req.RequestURI = ""
	req.URL, err = url.Parse(fmt.Sprintf("%s%s", server.URL, req.URL.Path))
	assert.Nil(t, err)
	client := http.DefaultClient
	resp, err := client.Do(req)
	fmt.Println(err)
	assert.Nil(t, err)
	fmt.Println(resp.StatusCode)
}
