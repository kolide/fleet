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

type handlerInfo struct {
	title       string
	path        string
	description string
	method      string
}

var hinfos = map[string]handlerInfo{
	"Login": {
		path:        "/api/v1/kolide/login",
		method:      "POST",
		description: "Logs a user in.",
		title:       "Login",
	},
	"Logout": {
		path:        "/api/v1/kolide/logout",
		method:      "POST",
		description: "Logs a user out.",
		title:       "Logout",
	},
	"Forgot Password": {
		path:        "/api/v1/kolide/forgot_password",
		method:      "POST",
		description: "Sends a password reset email.",
		title:       "Forgot Password",
	},
	"Me": {
		path:        "/api/v1/kolide/me",
		method:      "GET",
		description: "Returns info about the currently logged in user",
		title:       "Current Session User",
	},
	"GetUser": {
		path:        "/api/v1/kolide/users/1",
		method:      "GET",
		description: "Return a single user by ID",
		title:       "Get User",
	},
	"ListUsers": {
		path:        "/api/v1/kolide/users",
		method:      "GET",
		description: "Return a list of users",
		title:       "List Users",
	},
	"GetSessionInfo": {
		path:        "/api/v1/kolide/sessions/1",
		method:      "GET",
		description: "Return session info",
		title:       "Get Session Info",
	},
	"DeleteSession": {
		path:        "/api/v1/kolide/sessions/1",
		method:      "DELETE",
		description: "Delete session by ID",
		title:       "Delete Session",
	},
	"ListInvites": {
		path:        "/api/v1/kolide/invites",
		method:      "GET",
		description: "List invited users",
		title:       "List Invites",
	},
	"DeleteInvite": {
		path:        "/api/v1/kolide/invites/1",
		method:      "DELETE",
		description: "Delete invite by ID.",
		title:       "Delete Invite",
	},
	"ListHosts": {
		path:        "/api/v1/kolide/hosts",
		method:      "GET",
		description: "List hosts",
		title:       "List Hosts",
	},
	"GetHost": {
		path:        "/api/v1/kolide/hosts/1",
		method:      "GET",
		description: "Get a host by ID",
		title:       "Get Host",
	},
	"DeleteHost": {
		path:        "/api/v1/kolide/hosts/1",
		method:      "DELETE",
		description: "Delete a host by ID",
		title:       "Delete Host",
	},
}

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

	var docTests = []struct {
		info         handlerInfo
		body         io.Reader
		sessionToken string
	}{
		{
			info: hinfos["Login"],
			body: newBody([]byte(`{
        		"username": "admin1",
        		"password": "foobar"
			}`)),
		},
		{
			info: hinfos["Login"],
			body: newBody([]byte(`{
        		"username": "nobody",
        		"password": "foobar"
			}`)),
		},
		{
			info:         hinfos["Logout"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info: hinfos["Forgot Password"],
			body: newBody([]byte(`{"email":"nobody@kolide.co"}`)),
		},
		/*{
			path:   "/api/v1/kolide/forgot_password",
			method: "POST",
			body:   newBody([]byte(`{"email":"admin1@example.com"}`)),
		},
		TODO fails for now because mailService is nil
		*/
		{
			info:         hinfos["Me"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["ListUsers"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["GetUser"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["GetSessionInfo"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["DeleteSession"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["ListInvites"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["DeleteInvite"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["ListHosts"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["GetHost"],
			sessionToken: loggedInSessionToken(t, handler),
		},
		{
			info:         hinfos["DeleteHost"],
			sessionToken: loggedInSessionToken(t, handler),
		},
	}

	client := http.DefaultClient
	for _, tt := range docTests {
		t.Run(tt.info.title, func(t *testing.T) {
			req := httptest.NewRequest(tt.info.method, tt.info.path, tt.body)
			if tt.sessionToken != "" {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.sessionToken))
			}
			req.Header.Add("X-Test2Doc-Description", tt.info.description)
			req.Header.Add("X-Test2Doc-Title", tt.info.title)
			req.RequestURI = ""
			req.URL, err = url.Parse(fmt.Sprintf("%s%s", server.URL, req.URL.Path))
			require.Nil(t, err)
			_, err = client.Do(req)
			require.Nil(t, err)
		})
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
