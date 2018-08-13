package queue
import (
	"encoding/json"
	"context"
)

type InmemQueue struct {
	Memory         []json.RawMessage
}


func NewInmemQueue() (*InmemQueue, error) {
	return &InmemQueue{}, nil 
}


func (iq *InmemQueue) Messages(ctx context.Context, logs []json.RawMessage) error {
	iq.Memory = append(iq.Memory, logs...)
	return nil 
}
