package kolide

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/net/context"
)

type OptionStore interface {
	SaveOptions(opts []Option) error
	Options() ([]Option, error)
	Option(id uint) (*Option, error)
	OptionByName(name string) (*Option, error)
}

type OptionService interface {
	GetOptions(ctx context.Context) ([]Option, error)
	ModifyOptions(ctx context.Context, req OptionRequest) ([]Option, error)
}

const (
	ReadOnly    = true
	NotReadOnly = !ReadOnly
)

// OptionType defines the type of the value assigned to an option
type OptionType int

const (
	OptionTypeString OptionType = iota
	OptionTypeInt
	OptionTypeFlag
)

func (ot OptionType) String() string {
	switch ot {
	case OptionTypeString:
		return "string"
	case OptionTypeInt:
		return "int"
	case OptionTypeFlag:
		return "flag"
	}
	panic("stringer not implemented for OptionType")
}

// Option represents a possible osquery confguration option
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type Option struct {
	// ID unique identifier for option assigned by the dbms
	ID uint
	// Name the name of the option which must be unique
	Name string
	// Type of value expected for the option, db only
	Type OptionType
	// RawValue is string representation of option value, may be nil to
	// indicate the option is not set
	Value *string `db:"value"`
	// ReadOnly if true, this option is required for Kolide to function
	// properly and cannot be modified by the end user
	ReadOnly bool `db:"read_only"`
}

// OptionRequest contains options that are passed to modify options requests.
type OptionRequest struct {
	Options []Option `json:"options"`
}

type transform struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ReadOnly bool        `json:"read_only"`
}

func (opt *Option) MarshalJSON() ([]byte, error) {
	var val interface{}
	if opt.Value == nil {
		val = opt.Value
	} else {
		switch opt.Type {
		case OptionTypeString:
			val = *opt.Value
		case OptionTypeInt:
			i, err := strconv.ParseInt(*opt.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			val = i
		case OptionTypeFlag:
			b, err := strconv.ParseBool(*opt.Value)
			if err != nil {
				return nil, err
			}
			val = b
		default:
			panic("unimplemented option type")
		}
	}
	transform := &transform{
		opt.ID,
		opt.Name,
		opt.Type.String(),
		val,
		opt.ReadOnly,
	}
	return json.Marshal(transform)

}

func (opt *Option) UnmarshalJSON(b []byte) error {
	transform := &transform{}
	err := json.Unmarshal(b, transform)
	if err != nil {
		return err
	}
	opt.ID = transform.ID
	opt.Name = transform.Name
	opt.ReadOnly = transform.ReadOnly
	switch transform.Type {
	case "flag":
		opt.Type = OptionTypeFlag
	case "int":
		opt.Type = OptionTypeInt
	case "string":
		opt.Type = OptionTypeString
	default:
		return fmt.Errorf("option type '%s' invalid", transform.Type)
	}
	if transform.Value == nil {
		opt.Value = nil
		return nil
	}

	opt.Value = new(string)
	switch opt.Type {
	case OptionTypeFlag:
		v, ok := transform.Value.(bool)
		if !ok {
			return fmt.Errorf("option value type mismatch")
		}
		*opt.Value = strconv.FormatBool(v)
	case OptionTypeInt:
		v, ok := transform.Value.(float64)
		if !ok {
			return fmt.Errorf("option value type mismatch")
		}
		*opt.Value = strconv.Itoa(int(v))
	case OptionTypeString:
		v, ok := transform.Value.(string)
		if !ok {
			return fmt.Errorf("option value type mismatch")
		}
		*opt.Value = v
	}
	return nil
}
