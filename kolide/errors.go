package kolide

// An Error represents a Kolide error
type Error interface {
	error
	Code() ErrorCode
}

// An ErrorCode is a possible Kolide error
type ErrorCode int

const (
	// OK is returned on success.
	OK ErrorCode = iota

	// Unauthenticated indicates that credentials are missing or invalid
	Unauthenticated

	// Unauthorized indicates that the requester does not have permission
	// to make the request.
	Unauthorized

	// InvalidInput indicates that the request is malformed
	// or missing required arguments.
	InvalidInput

	// NotFound indicates that the requested resource cannot be found.
	NotFound

	// AlreadyExists indicates that a resource was not created
	// because it already exists.
	AlreadyExists

	// Unavailable indicates that the server is busy or
	// another temporary error has occured.
	// Unavailable errors are retriablie.
	Unavailable

	// Unknown indicates that an error which is not expected has occured.
	Unknown = 500
)
