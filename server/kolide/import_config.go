package kolide

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

type ImportConfigService interface {
	// ImportOsqueryConfiguration create packs, queries, options etc based on imported
	// osquery configuration
	ImportConfig(ctx context.Context, cfg *ImportConfig) (*ImportConfigResponse, error)
}

// ImportSection is used to categorize information associated with the import
// of a particular section of an imported osquery configuration file
type ImportSection string

const (
	OptionsSection ImportSection = "options"
	PacksSection                 = "packs"
	QueriesSection               = "queries"
)

// WarningType is used to group associated warnings for options, packs etc
// when importing on osquery configuration file
type WarningType string

const (
	PackDuplicate    WarningType = "duplicate_pack"
	OptionAlreadySet             = "option_already_set"
	OptionReadonly               = "option_readonly"
	OptionUnknown                = "option_unknown"
)

// ImportStatus contains information pertaining to the import of a section
// of an osquery configuration file
type ImportStatus struct {
	// Title human readable name of the section of the import file that this
	// status pertains to.
	Title string `json:"title"`
	// ImportCount count of items successfully imported
	ImportCount int `json:"import_count"`
	// SkipCount count of items that are skipped.  The reasons for the omissions
	// can be found in Warnings
	SkipCount int `json:"skip_count"`
	// Warnings groups catagories of warnings with one or more detail messages
	Warnings map[WarningType][]string `json:"warnings"`
	// Messages contains an entry for each import attempt
	Messages []string `json:"messages"`
}

// Warning is used to add a warning message to ImportStatus
func (is *ImportStatus) Warning(warnType WarningType, fmtMsg string, fmtArgs ...interface{}) {
	is.Warnings[warnType] = append(is.Warnings[warnType], fmt.Sprintf(fmtMsg, fmtArgs...))
}

// Message is used to add a general message to ImportStatus, usually indicating
// what was changed in a successful import
func (is *ImportStatus) Message(fmtMsg string, args ...interface{}) {
	is.Messages = append(is.Messages, fmt.Sprintf(fmtMsg, args...))
}

// ImportConfigResponse contains information about the import of an osquery
// configuration file
type ImportConfigResponse struct {
	ImportStatusBySection map[ImportSection]*ImportStatus `json:"import_status"`
}

func NewImportConfigResponse() *ImportConfigResponse {
	return &ImportConfigResponse{
		ImportStatusBySection: make(map[ImportSection]*ImportStatus),
	}
}

// Status returns a structure that contains information about the import
// of a particular section of an osquery configuration file.
func (ic *ImportConfigResponse) Status(section ImportSection) (status *ImportStatus) {
	var ok bool
	if status, ok = ic.ImportStatusBySection[section]; !ok {
		status = new(ImportStatus)
		status.Title = strings.Title(string(section))
		status.Warnings = make(map[WarningType][]string)
		ic.ImportStatusBySection[section] = status
	}
	return status
}

const (
	GlobPacks = "*"
)

// QueryDetails represents the query objects used in the packs and the
// schedule section of an osquery configuration
type QueryDetails struct {
	Query    string `json:"query"`
	Interval int    `json:"interval"`
	// Optional fields
	Removed  *bool   `json:"removed"`
	Platform *string `json:"platform"`
	Version  *string `json:"version"`
	Shard    *int    `json:"shard"`
	Snapshot *bool   `json:"snapshot"`
}

// PackDetails represents the "packs" section of an osquery configuration
// file
type PackDetails struct {
	Queries   []QueryDetails `json:"queries"`
	Shard     *int           `json:"shard"`
	Version   *string        `json:"version"`
	Platform  *string        `json:"platform"`
	Discovery []string       `json:"discovery"`
}

// YARAConfig yara configuration maps keys to lists of files
// See https://osquery.readthedocs.io/en/stable/deployment/yara/
type YARAConfig struct {
	Signatures map[string][]string `json:"signatures"`
	FilePaths  map[string][]string `json:"file_paths"`
}

// Decorator section of osquery config each section contains rows of decorator
// queries
type DecoratorConfig struct {
	Load   []string `json:"load"`
	Always []string `json:"always"`
	/*
		Interval maps a string representation of a numeric interval to a set
		of decorator queries.
			{
				"interval": {
					"3600": [
						"SELECT total_seconds FROM uptime;"
					]
				}
			}
	*/
	Interval map[string][]string `json:"interval"`
}

type OptionNameToValueMap map[string]interface{}
type QueryNameToQueryDetailsMap map[string]QueryDetails
type PackNameMap map[string]interface{}
type FIMCategoryToPaths map[string][]string
type PackNameToPackDetails map[string]PackDetails

// ImportConfig is a representation of an Osquery configuration. Osquery
// documentation has further details.
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type ImportConfig struct {
	// Options is a map of option name to a value which can be an int,
	// bool, or string
	Options OptionNameToValueMap `json:"options"`
	// Schedule is a map of query names to details
	Schedule QueryNameToQueryDetailsMap `json:"schedule"`
	// Packs is a map of pack names to either PackDetails, or a string
	// containing a file path with a pack config. If a string, we expect
	// PackDetails to be stored in ExternalPacks
	Packs PackNameMap `json:"packs"`
	// FileIntegrityMonitoring file integrity monitoring information
	// See https://osquery.readthedocs.io/en/stable/deployment/file-integrity-monitoring/
	FileIntegrityMonitoring FIMCategoryToPaths `json:"file_paths"`
	// YARA configuration
	YARA       *YARAConfig      `json:"yara"`
	Decorators *DecoratorConfig `json:"decorators"`
	// ExternalPacks are packs referenced when an item in the Packs map references
	// an external file.  The PackName here must match the PackName in the Packs map
	ExternalPacks PackNameToPackDetails `json:"-"`
	// GlobPackNames lists pack names that are globbed
	GlobPackNames []string
}
