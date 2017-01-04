package kolide

import (
	"golang.org/x/net/context"
)

type ImportConfigService interface {
	// ImportOsqueryConfiguration create packs, queries, options etc based on imported
	// osquery configuration
	ImportConfig(ctx context.Context, cfg *ImportConfig) (*ImportConfigResponse, error)
}

type ImportConfigResponse struct {
}

// type ImportConfigRequest struct {
// 	Config ImportConfig `json:"config_import"`
// }

type QueryName string
type OptionName string

const (
	GlobPacks = "*"
)

type FIMCategory string
type IntervalValue string

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
	Version   *string        `json:"shard"`
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
	Interval map[IntervalValue][]string `json:"interval"`
}

// ImportConfig is a representation of an Osquery configuration. Osquery
// documentation has further details.
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type ImportConfig struct {
	// Options is a map of option name to a value which can be an int,
	// bool, or string
	Options map[OptionName]interface{} `json:"options"`
	// Schedule is a map of query names to details
	Schedule map[QueryName]QueryDetails `json:"schedule"`
	// Packs is a map of pack names to either PackDetails, or a string
	// containing a file path with a pack config. If a string, we expect
	// PackDetails to be stored in ExternalPacks
	Packs map[string]interface{} `json:"packs"`
	// FileIntegrityMonitoring file integrity monitoring information
	// See https://osquery.readthedocs.io/en/stable/deployment/file-integrity-monitoring/
	FileIntegrityMonitoring map[FIMCategory][]string `json:"file_paths"`
	// YARA configuration
	YARA       YARAConfig      `json:"yara"`
	Decorators DecoratorConfig `json:"decorators"`
	// ExternalPacks are packs referenced when an item in the Packs map references
	// an external file.  The PackName here must match the PackName in the Packs map
	ExternalPacks map[string]PackDetails `json:"-"`
	// GlobPackNames lists pack names that are globbed
	GlobPackNames []string
}
