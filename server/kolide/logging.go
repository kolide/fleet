package kolide

import "encoding/json"

// OsqueryLogger defines an interface for loggers that can write osquery JSON
// to various output sources.
type OsqueryLogger interface {
	// Write writes the JSON log entries to the appropriate destination,
	// returning any errors that occurred.
	Write(logs []json.RawMessage) error
}
