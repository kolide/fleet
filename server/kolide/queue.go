package kolide

import (
	"context"
	"encoding/json"
)

// QueueService defines functions for sending query results over to external
// system. It is implemented by structs in package queue.
type QueueService interface {
	Messages(ctx context.Context, logs []json.RawMessage) error
}
