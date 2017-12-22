package service

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/contexts/token"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

const bcryptCost = 6

func TestAuthenticate(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)
	svc, err := newTestService(ds, nil)
	require.Nil(t, err)
	users := createTestUsers(t, ds)

	var loginTests = []struct {
		username string
		password string
		user     kolide.User
		wantErr  error
	}{
		{
			user:     users["admin1"],
			username: testUsers["admin1"].Username,
			password: testUsers["admin1"].PlaintextPassword,
		},
		{
			user:     users["user1"],
			username: testUsers["user1"].Email,
			password: testUsers["user1"].PlaintextPassword,
		},
	}

	for _, tt := range loginTests {
		t.Run(tt.username, func(st *testing.T) {
			user := tt.user
			ctx := context.Background()
			loggedIn, token, err := svc.Login(ctx, tt.username, tt.password)
			require.Nil(st, err, "login unsuccesful")
			assert.Equal(st, user.ID, loggedIn.ID)
			assert.NotEmpty(st, token)

			sessions, err := svc.GetInfoAboutSessionsForUser(ctx, user.ID)
			require.Nil(st, err)
			require.Len(st, sessions, 1, "user should have one session")
			session := sessions[0]
			assert.Equal(st, user.ID, session.UserID)
			assert.WithinDuration(st, time.Now(), session.AccessedAt, 3*time.Second,
				"access time should be set with current time at session creation")
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	jwtKey := ""
	tokenString, err := generateJWT("4", jwtKey)
	require.Nil(t, err)

	svc := authViewerService{}
	viewer, err := authViewer(
		context.Background(),
		jwtKey,
		token.Token(tokenString),
		svc,
	)
	require.Nil(t, err)
	require.NotNil(t, viewer)
}

type authViewerService struct {
	kolide.Service
}

func (authViewerService) GetSessionByKey(ctx context.Context, key string) (*kolide.Session, error) {
	return &kolide.Session{}, nil
}

func (authViewerService) User(ctx context.Context, uid uint) (*kolide.User, error) {
	return &kolide.User{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func TestDecodeGetInfoAboutSessionRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/sessions/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeGetInfoAboutSessionRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(getInfoAboutSessionRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/sessions/1", nil),
	)
}

func TestDecodeGetInfoAboutSessionsForUserRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/user/{id}/sessions", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeGetInfoAboutSessionsForUserRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(getInfoAboutSessionsForUserRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/users/1/sessions", nil),
	)
}

func TestDecodeDeleteSessionRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/sessions/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeDeleteSessionRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(deleteSessionRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("DELETE")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("DELETE", "/api/v1/kolide/sessions/1", nil),
	)
}

func TestDecodeDeleteSessionsForUserRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/user/{id}/sessions", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeDeleteSessionsForUserRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(deleteSessionsForUserRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("DELETE")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("DELETE", "/api/v1/kolide/users/1/sessions", nil),
	)
}

func TestDecodeLoginRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/login", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeLoginRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(loginRequest)
		assert.Equal(t, "foo", params.Username)
		assert.Equal(t, "bar", params.Password)
	}).Methods("POST")
	t.Run("lowercase username", func(t *testing.T) {
		var body bytes.Buffer
		body.Write([]byte(`{
        "username": "foo",
        "password": "bar"
    }`))

		router.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("POST", "/api/v1/kolide/login", &body),
		)
	})
	t.Run("uppercase username", func(t *testing.T) {
		var body bytes.Buffer
		body.Write([]byte(`{
        "username": "Foo",
        "password": "bar"
    }`))

		router.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("POST", "/api/v1/kolide/login", &body),
		)
	})

}
