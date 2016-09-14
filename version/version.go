// Package version provides utilities for displaying version information
package version

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// info vars
// set at runtime with ldflags
// example:
// go build -ldflags "-X github.com/kolide/kolide-ose/version.version=1.0.0-development"
var (
	version   = "unknown"
	branch    = "unknown"
	revision  = "unknown"
	goVersion = "unknown"
	buildDate = "unknown"
	buildUser = "unknown"
)

// Info holds version and build info about Kolide
type Info struct {
	Version   string `json:"version"`
	Branch    string `json:"branch"`
	Revision  string `json:"revision"`
	GoVersion string `json:"go_version"`
	BuildDate string `json:"build_date"`
	BuildUser string `json:"build_user"`
}

// Version returns a struct with the current version information
func Version() Info {
	return Info{
		Version:   version,
		Branch:    branch,
		Revision:  revision,
		GoVersion: goVersion,
		BuildDate: buildDate,
		BuildUser: buildUser,
	}
}

// Print outputs the app name and version string
func Print() {
	v := Version()
	fmt.Printf("kolide version %s\n", v.Version)
}

// PrintFull outputs the app name and detailed version information
func PrintFull() {
	v := Version()
	fmt.Printf("kolide - version %s\n", v.Version)
	fmt.Printf("  branch: \t%s\n", v.Branch)
	fmt.Printf("  revision: \t%s\n", v.Revision)
	fmt.Printf("  go version: \t%s\n", v.GoVersion)
}

func Handler() http.Handler {
	v := Version()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(v)

	})
}
