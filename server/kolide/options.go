package kolide

import "golang.org/x/net/context"

type OptionStore interface {
	NewOption(name string, optType OptionType, kolideRequires bool) (*Option, error)
	Options() ([]Option, error)
	SetOptionValue(optID uint, value interface{}) (*OptionValue, error)
	OptionValues() ([]OptionValue, error)
	DeleteOptionValue(optID uint) error
}

type OptionService interface {
	Options(ctx context.Context) (*OptionsPayload, error)
	OptionValues(ctx context.Context) (*OptionsPayload, error)
	ModifyOptionValues(ctx context.Context, p *OptionValuesPayload) (*OptionValuesPayload, error)
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
	// RequiredForKolide if true, this option is required for Kolide to function
	// properly and cannot be modified by the end user
	RequiredForKolide bool `json:"required_for_kolide" db:"required_for_kolide"`
}

// OptionTypePayload defines the json representation of options exposed via the
// API
type OptionsPayload struct {
	Options []Option `json:"options"`
}

// OptionValue represents a value used to configure Osquery options
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type OptionValue struct {
	UpdateCreateTimestamps
	ID       uint   `json:id`
	OptionID uint   `json:"option_id" db:"option_id"`
	Value    string `json:"-"`
	// OptionValue contains the value of the option set by the end user, it can
	// be int, string or bool
	OptionValue interface{} `json:"value" db:"-"`
}

// OptionValuesPayload defines the json representation of option values that
// will be exposed via the API
type OptionValuesPayload struct {
	OptionValues []OptionValue `json:"option_values"`
}
