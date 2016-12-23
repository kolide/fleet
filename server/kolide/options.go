package kolide

type OptionStore interface {
	SaveOption(opt Option) error
	Options() ([]Option, error)
	Option(id uint) (*Option, error)
	OptionByName(name string) (*Option, error)
}

type OptionService interface {
	// Options(ctx context.Context) (*OptionsPayload, error)
	// OptionValues(ctx context.Context) (*OptionsPayload, error)
	// ModifyOptionValues(ctx context.Context, p *OptionValuesPayload) (*OptionValuesPayload, error)
}

// OptionType defines the type of the value assigned to an option
type OptionType int

const (
	OptionTypeString OptionType = iota
	OptionTypeInt
	OptionTypeFlag
)

// Option represents a possible osquery confguration option
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type Option struct {
	// ID unique identifier for option assigned by the dbms
	ID uint `json:"id"`
	// Name the name of the option which must be unique
	Name string `json:"name"`
	// Type of value expected for the option, db only
	Type OptionType `json:"-"`
	// TypeName maps a front end friendly tag to Type
	TypeName string `json:"type" db:"-"`
	// RawValue is string representation of option value, may be nil to
	// indicate the option is not set
	RawValue *string `json:"-" db:"value"`
	// Value is the option value which can be bool, int, or string
	Value interface{} `json:"value" db:"-"`
	// ReadOnly if true, this option is required for Kolide to function
	// properly and cannot be modified by the end user
	ReadOnly bool `json:"read_only" db:"read_only"`
}

// OptionTypePayload defines the json representation of options exposed via the
// API
type OptionsPayload struct {
	Options []Option `json:"options"`
}

const (
	ReadOnly    = true
	NotReadOnly = !ReadOnly
)
