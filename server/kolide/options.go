package kolide

import "golang.org/x/net/context"

type OptionStore interface {
	NewOption(name string, optType OptionType, kolideRequires bool) (*Option, error)
	Options() ([]Option, error)
	SetOptionValues([]OptionValue) ([]OptionValue, error)
	OptionValues() ([]OptionValue, error)
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

const (
	RequiredByKolide = true
	OptionChangable  = false
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
	ID       uint   `json:"id"`
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

// InitializeOptions creates osquery option values
func InitializeOptions(ds Datastore) error {
	if _, err := ds.NewOption("aws_access_key_id", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_firehose_period", OptionTypeInt, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_firehose_stream", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_kinesis_period", OptionTypeInt, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_kinesis_random_partition_key", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_kinesis_stream", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_profile_name", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_region", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_secret_access_key", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_sts_arn_role", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_sts_region", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_sts_session_name", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("aws_sts_timeout", OptionTypeInt, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("buffered_log_max", OptionTypeInt, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("decorations_top_level", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_caching", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_database", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_decorators", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_distributed", OptionTypeFlag, true); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_events", OptionTypeFlag, true); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_kernel", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_logging", OptionTypeFlag, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("disable_tables", OptionTypeString, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("distributed_interval", OptionTypeInt, false); err != nil {
		return err
	}
	if _, err := ds.NewOption("distributed_plugin", OptionTypeString, true); err != nil {
		return err
	}
	return nil

}
