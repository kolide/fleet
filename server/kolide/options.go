package kolide

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/net/context"
)

// OptionStore interface describes methods to access datastore
type OptionStore interface {
	SaveOptions(opts []Option) error
	Options() ([]Option, error)
	Option(id uint) (*Option, error)
	OptionByName(name string) (*Option, error)
}

// OptionsService interface describes methods that operate on osquery options
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

// Used to marshal OptionType to human readable strings used in JSON payloads
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
	Value interface{} `db:"value"`
	// ReadOnly if true, this option is required for Kolide to function
	// properly and cannot be modified by the end user
	ReadOnly bool `db:"read_only"`
}

// ValueAsString is a convenience method that converts the Value field of an
// option to a string.
func (opt *Option) ValueAsString() (val string, isDefined bool) {
	if opt.Value == nil {
		return "", false
	}
	switch opt.Value.(type) {
	case []uint8:
		return string(opt.Value.([]uint8)), true
	case *string:
		return *(opt.Value.(*string)), true
	case string:
		return opt.Value.(string), true
	}
	panic("unknown option type")
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
		val = nil
	} else {
		strPtr, ok := opt.Value.(*string)
		if !ok {
			return nil, fmt.Errorf("option value is not expected type")
		}
		if strPtr == nil {
			val = strPtr
		} else {
			switch opt.Type {
			case OptionTypeString:
				val = *strPtr
			case OptionTypeInt:
				i, err := strconv.ParseInt(*strPtr, 10, 64)
				if err != nil {
					return nil, err
				}
				val = i
			case OptionTypeFlag:
				b, err := strconv.ParseBool(*strPtr)
				if err != nil {
					return nil, err
				}
				val = b
			default:
				panic("unimplemented option type")
			}
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
	opt.Value = transform.Value
	return nil
}
