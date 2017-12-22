package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

////////////////////////////////////////////////////////////////////////////////
// Service Errors
////////////////////////////////////////////////////////////////////////////////

type invalidArgumentError []invalidArgument
type invalidArgument struct {
	name   string
	reason string
}

// newInvalidArgumentError returns a invalidArgumentError with at least
// one error.
func newInvalidArgumentError(name, reason string) *invalidArgumentError {
	var invalid invalidArgumentError
	invalid = append(invalid, invalidArgument{
		name:   name,
		reason: reason,
	})
	return &invalid
}

func (e *invalidArgumentError) Append(name, reason string) {
	*e = append(*e, invalidArgument{
		name:   name,
		reason: reason,
	})
}
func (e *invalidArgumentError) Appendf(name, reasonFmt string, args ...interface{}) {
	*e = append(*e, invalidArgument{
		name:   name,
		reason: fmt.Sprintf(reasonFmt, args...),
	})
}

func (e *invalidArgumentError) HasErrors() bool {
	return len(*e) != 0
}

// invalidArgumentError is returned when one or more arguments are invalid.
func (e invalidArgumentError) Error() string {
	switch len(e) {
	case 0:
		return "validation failed"
	case 1:
		return fmt.Sprintf("validation failed: %s %s", e[0].name, e[0].reason)
	default:
		return fmt.Sprintf("validation failed: %s %s and %d other errors", e[0].name, e[0].reason,
			len(e))
	}
}

func (e invalidArgumentError) Invalid() []map[string]string {
	var invalid []map[string]string
	for _, i := range e {
		invalid = append(invalid, map[string]string{"name": i.name, "reason": i.reason})
	}
	return invalid
}

// authentication error
type authError struct {
	reason string
	// client reason is used to provide
	// a different error message to the client
	// when security is a concern
	clientReason string
}

func (e authError) Error() string {
	return e.reason
}

func (e authError) AuthError() string {
	if e.clientReason != "" {
		return e.clientReason
	}
	return "username or email and password do not match"
}

// permissionError, set when user is authenticated, but not allowed to perform action
type permissionError struct {
	message string
	badArgs []invalidArgument
}

func newPermissionError(name, reason string) permissionError {
	return permissionError{
		badArgs: []invalidArgument{
			invalidArgument{
				name:   name,
				reason: reason,
			},
		},
	}
}

func (e permissionError) Error() string {
	switch len(e.badArgs) {
	case 0:
	case 1:
		e.message = fmt.Sprintf("unauthorized: %s",
			e.badArgs[0].reason,
		)
	default:
		e.message = fmt.Sprintf("unauthorized: %s and %d other errors",
			e.badArgs[0].reason,
			len(e.badArgs),
		)
	}
	if e.message == "" {
		return "unauthorized"
	}
	return e.message
}

func (e permissionError) PermissionError() []map[string]string {
	var forbidden []map[string]string
	if len(e.badArgs) == 0 {
		forbidden = append(forbidden, map[string]string{"reason": e.Error()})
		return forbidden
	}
	for _, arg := range e.badArgs {
		forbidden = append(forbidden, map[string]string{
			"name":   arg.name,
			"reason": arg.reason,
		})
	}
	return forbidden

}

////////////////////////////////////////////////////////////////////////////////
// Transport Errors
////////////////////////////////////////////////////////////////////////////////

// erroer interface is implemented by response structs to encode business logic errors
type errorer interface {
	error() error
}

type jsonError struct {
	Message string              `json:"message"`
	Errors  []map[string]string `json:"errors,omitempty"`
}

// use baseError to encode an jsonError.Errors field with an error that has
// a generic "name" field. The frontend client always expects errors in a
// []map[string]string format.
func baseError(err string) []map[string]string {
	return []map[string]string{map[string]string{
		"name":   "base",
		"reason": err},
	}
}

// same as baseError, but replaces "base" with different name.
func namedError(name string, err string) []map[string]string {
	return []map[string]string{map[string]string{
		"name":   name,
		"reason": err},
	}
}

// encode error and status header to the client
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	type validationError interface {
		error
		Invalid() []map[string]string
	}
	if e, ok := err.(validationError); ok {
		ve := jsonError{
			Message: "Validation Failed",
			Errors:  e.Invalid(),
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		enc.Encode(ve)
		return
	}

	type authenticationError interface {
		error
		AuthError() string
	}
	if e, ok := err.(authenticationError); ok {
		ae := jsonError{
			Message: "Authentication Failed",
			Errors:  baseError(e.AuthError()),
		}
		w.WriteHeader(http.StatusUnauthorized)
		enc.Encode(ae)
		return
	}

	type permissionError interface {
		PermissionError() []map[string]string
	}
	if e, ok := err.(permissionError); ok {
		pe := jsonError{
			Message: "Permission Denied",
			Errors:  e.PermissionError(),
		}
		w.WriteHeader(http.StatusForbidden)
		enc.Encode(pe)
		return
	}

	type mailError interface {
		MailError() []map[string]string
	}
	if e, ok := err.(mailError); ok {
		me := jsonError{
			Message: "Mail Error",
			Errors:  e.MailError(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(me)
		return
	}

	type osqueryError interface {
		error
		NodeInvalid() bool
	}
	if e, ok := err.(osqueryError); ok {
		// osquery expects to receive the node_invalid key when a TLS
		// request provides an invalid node_key for authentication. It
		// doesn't use the error message provided, but we provide this
		// for debugging purposes (and perhaps osquery will use this
		// error message in the future).

		errMap := map[string]interface{}{"error": e.Error()}
		if e.NodeInvalid() {
			w.WriteHeader(http.StatusUnauthorized)
			errMap["node_invalid"] = true
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		enc.Encode(errMap)
		return
	}

	type notFoundError interface {
		error
		IsNotFound() bool
	}
	if e, ok := err.(notFoundError); ok {
		je := jsonError{
			Message: "Resource Not Found",
			Errors:  baseError(e.Error()),
		}
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(je)
		return
	}

	type existsError interface {
		error
		IsExists() bool
	}
	if e, ok := err.(existsError); ok {
		je := jsonError{
			Message: "Resource Already Exists",
			Errors:  baseError(e.Error()),
		}
		w.WriteHeader(http.StatusConflict)
		enc.Encode(je)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	je := jsonError{
		Message: "Unknown Error",
		Errors:  baseError(err.Error()),
	}
	enc.Encode(je)
}
