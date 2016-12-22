package kolide

import (
	"fmt"
	"strconv"
)

type OptionStore interface {
	SaveOption(opt Option) (*Option, error)
	Options() ([]Option, error)
	Option(id uint) (*Option, error)
	// SetOptionValues([]OptionValue) ([]OptionValue, error)
	// OptionValues() ([]OptionValue, error)
	// GetOptionInt(name string) (int, error)
	// GetOptionFlag(name string) (bool, error)
	// GetOptionString(name string) (string, error)
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

func NewOption(name string, value interface{}, typ OptionType, readonly bool) (*Option, error) {
	//rawValue = new(string)
	result := &Option{
		Name:     name,
		Type:     typ,
		Value:    value,
		ReadOnly: readonly,
	}
	if value == nil {
		return result, nil
	}
	result.RawValue = new(string)
	switch typ {
	case OptionTypeFlag:
		boolVal, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("type mismatch")
		}
		*result.RawValue = strconv.FormatBool(boolVal)
	case OptionTypeString:
		ok := false
		if *result.RawValue, ok = value.(string); !ok {
			return nil, fmt.Errorf("type mismatch")
		}
	case OptionTypeInt:
		intVal, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("type mismatch")
		}
		*result.RawValue = strconv.FormatInt(int64(intVal), 10)
	}
	return result, nil
}

// OptionTypePayload defines the json representation of options exposed via the
// API
type OptionsPayload struct {
	Options []Option `json:"options"`
}

// InitializeOptions creates osquery option values
func CreateOptions(ds Datastore) error {
	// if _, err := ds.NewOption("aws_access_key_id", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_firehose_period", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_firehose_stream", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_kinesis_period", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_kinesis_random_partition_key", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_kinesis_stream", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_profile_name", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_region", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_secret_access_key", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_sts_arn_role", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_sts_region", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_sts_session_name", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("aws_sts_timeout", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("buffered_log_max", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("decorations_top_level", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_caching", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_database", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_decorators", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_distributed", OptionTypeFlag, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_events", OptionTypeFlag, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_kernel", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_logging", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("disable_tables", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("distributed_interval", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("distributed_plugin", OptionTypeString, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("distributed_tls_max_attempts", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("distributed_tls_read_endpoint", OptionTypeString, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("distributed_tls_write_endpoint", OptionTypeString, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("enable_foreign", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("enable_monitor", OptionTypeFlag, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("ephemeral", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("events_expiry", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("events_max", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("events_optimize", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("host_identifier", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_event_type", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_mode", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_path", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_plugin", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_secondary_status_only", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_syslog_facility", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_tls_compress", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_tls_endpoint", OptionTypeString, RequiredByKolide); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_tls_max", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("logger_tls_period", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("pack_delimiter", OptionTypeString, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("pack_refresh_interval", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("read_max", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("read_user_max", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("schedule_default_interval", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("schedule_splay_percent", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("schedule_timeout", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("utc", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("value_max", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("verbose", OptionTypeFlag, OptionChangable); err != nil {
	// 	return err
	// }
	// if _, err := ds.NewOption("worker_threads", OptionTypeInt, OptionChangable); err != nil {
	// 	return err
	// }
	return nil

}
