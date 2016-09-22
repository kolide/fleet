// Package token enables setting and reading
// authentication token contexts
package token

import (
	"context"
	"net/http"
	"strings"
)

type key int

const tokenKey key = 0

// NewContext returns a new context carrying the Authorization Bearer token.
func NewContext(ctx context.Context, r *http.Request) context.Context {
	headers := r.Header.Get("Authorization")
	headerParts := strings.Split(headers, " ")
	if len(headerParts) != 2 || strings.ToUpper(headerParts[0]) != "BEARER" {
		return ctx
	}
	return context.WithValue(ctx, tokenKey, headerParts[1])
}

// FromContext extracts the Authorization Bearer token if present.
func FromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	return token, ok
}
