package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/kolide-ose/kolide"
)

func loginT(svc kolide.Service, logger kitlog.Logger) http.HandlerFunc {
	// TODO: pass from config
	var (
		secret               = []byte("insecure-jwt-key")
		authTokenDuration    = time.Minute * 30
		refreshTokenDuration = time.Hour * 24 * 14
	)
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest struct {
			Username *string
			Password *string
		}
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			encodeResponse(ctx, w, getUserResponse{
				Err: err,
			})
			logger.Log("err", err)
			return
		}
		var username, password string
		{
			if loginRequest.Username != nil {
				username = *loginRequest.Username
			}
			if loginRequest.Password != nil {
				password = *loginRequest.Password
			}
		}
		// retrieve user or respond with error
		user, err := svc.Authenticate(ctx, username, password)
		switch err.(type) {
		case nil:
			logger.Log("msg", "authenticated", "user", username, "id", user.ID)
		case authError:
			encodeResponse(ctx, w, getUserResponse{
				Err: err,
			})
			logger.Log("err", err, "user", username)
			return
		default:
			encodeResponse(ctx, w, getUserResponse{
				Err: errors.New("unknown error, try again later"),
			})
			logger.Log("err", err, "user", username)
			return
		}

		// logged in, give out a token
		authToken, err := newToken(authTokenDuration, user.ID, secret)
		if err != nil {
			encodeResponse(ctx, w, getUserResponse{
				Err: errors.New("unknown error, try again later"),
			})
			logger.Log("err", err, "user", username)
			return
		}

		refreshToken, err := newToken(refreshTokenDuration, user.ID, secret)
		if err != nil {
			encodeResponse(ctx, w, getUserResponse{
				Err: errors.New("unknown error, try again later"),
			})
			logger.Log("err", err, "user", username)
			return
		}

		var jwtLoginResponse = struct {
			jwtToken
			getUserResponse
		}{
			jwtToken{
				AuthToken:    authToken,
				RefreshToken: refreshToken,
			},
			getUserResponse{
				ID:                 user.ID,
				Username:           user.Username,
				Name:               user.Name,
				Admin:              user.Admin,
				Enabled:            user.Enabled,
				NeedsPasswordReset: user.NeedsPasswordReset,
			},
		}

		encodeResponse(ctx, w, jwtLoginResponse)

	}
}

func refresh(svc kolide.Service, logger kitlog.Logger) http.HandlerFunc {
	var (
		secret            = []byte("insecure-jwt-key")
		authTokenDuration = time.Minute * 30
	)
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := tokenFromRequest(r, secret)
		if err != nil || !token.Valid {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			logger.Log("err", err)
			return
		}
		authToken, err := newToken(authTokenDuration, claims.UserID, secret)
		if err != nil {
			encodeResponse(ctx, w, getUserResponse{
				Err: errors.New("unknown error, try again later"),
			})
			logger.Log("err", err)
			return
		}

		encodeResponse(ctx, w, jwtToken{AuthToken: authToken})
	}
}

func authMiddlewareT(svc kolide.Service, logger kitlog.Logger, next http.Handler) http.Handler {
	var (
		secret = []byte("insecure-jwt-key")
	)
	logger = kitlog.NewContext(logger).With("method", "authMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			logger.Log("err", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func tokenFromRequest(r *http.Request, secret []byte) (*jwt.Token, *kolideClaims, error) {
	claims := &kolideClaims{}
	token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}

type jwtToken struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type kolideClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func newToken(duration time.Duration, userID uint, secret []byte) (string, error) {
	claims := kolideClaims{
		userID,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
