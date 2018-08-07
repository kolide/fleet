package kolide

import (
	"context"
	"encoding/json"
)

type QueueService interface {
	Messages(ctx context.Context, logs []json.RawMessage) error
}
