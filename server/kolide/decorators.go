package kolide

import (
	"errors"

	"golang.org/x/net/context"
)

// DecoratorStore methods to manipulate decorator queries.
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type DecoratorStore interface {
	// NewDecorator creates a decorator query.
	NewDecorator(decorator *Decorator) (*Decorator, error)
	// DeleteDecorator removes a decorator query.
	DeleteDecorator(id uint) error
	// Decorator retrieves a decorator query with supplied ID.
	Decorator(id uint) (*Decorator, error)
	// ListDecorators returns all decorator queries.
	ListDecorators() ([]*Decorator, error)
	// SaveDecorator updates an existing decorator
	SaveDecorator(dec *Decorator) error
}

// DecoratorService exposes decorators data so it can be manipulated by
// end users
type DecoratorService interface {
	// ListDecorators returns decorators
	ListDecorators(ctx context.Context) ([]*Decorator, error)
	// DeleteDecorator removes an existing decorator if it is not built-in
	DeleteDecorator(ctx context.Context, id uint) error
	// NewDecorator creates a new decorator
	NewDecorator(ctx context.Context, payload DecoratorPayload) (*Decorator, error)
	// ModifyDecorator updates an existing decorator
	ModifyDecorator(ctx context.Context, payload DecoratorPayload) (*Decorator, error)
}

// DecoratorType refers to the allowable types of decorator queries.
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type DecoratorType int

const (
	DecoratorLoad DecoratorType = iota
	DecoratorAlways
	DecoratorInterval
	DecoratorUndefined

	DecoratorLoadName     = `"load"`
	DecoratorAlwaysName   = `"always"`
	DecoratorIntervalName = `"interval"`
)

func (dt DecoratorType) String() string {
	switch dt {
	case DecoratorLoad:
		return DecoratorLoadName
	case DecoratorAlways:
		return DecoratorAlwaysName
	case DecoratorInterval:
		return DecoratorIntervalName
	default:
		return ""
	}
}

func (dt *DecoratorType) MarshalJSON() ([]byte, error) {
	name := dt.String()
	if name == "" {
		return nil, errors.New("Invalid decorator type")
	}
	return []byte(name), nil
}

func (dt *DecoratorType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case DecoratorLoadName:
		*dt = DecoratorLoad
	case DecoratorAlwaysName:
		*dt = DecoratorAlways
	case DecoratorIntervalName:
		*dt = DecoratorInterval
	default:
		return errors.New("Invalid load type")
	}
	return nil
}

var decNameToType map[string]DecoratorType

func DecoratorTypeFromName(name string) (DecoratorType, error) {
	if decNameToType == nil {
		decNameToType = map[string]DecoratorType{
			"load":     DecoratorLoad,
			"always":   DecoratorAlways,
			"interval": DecoratorInterval,
		}
	}
	if typ, ok := decNameToType[name]; ok {
		return typ, nil
	}
	return DecoratorUndefined, errors.New("undefined decorator type")
}

// Decorator contains information about a decorator query.
type Decorator struct {
	UpdateCreateTimestamps
	ID   uint          `json:"id"`
	Type DecoratorType `json:"type"`
	// Interval note this is only pertainent for DecoratorInterval type.
	Interval uint   `json:"interval"`
	Query    string `json:"query"`
	// BuiltIn decorators are loaded in migrations and may not be changed
	BuiltIn bool `json:"built_in" db:"built_in"`
}

type DecoratorPayload struct {
	ID            uint    `json:"id"`
	DecoratorType *string `json:"decorator_type"`
	Interval      *uint   `json:"interval"`
	Query         *string `json:"query"`
}
