package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"

	yaml "gopkg.in/yaml.v2"
)

type Settings struct {
	AllowedLicenses map[string]interface{} `yaml:"allowed_licenses"`
	Exceptions      map[string]interface{} `yaml:"exceptions"`
}

type Dependency struct {
	Name       string
	License    string
	Repository string
}

// getJavascriptDependencies retrieves the licensing metadata for javascript
// dependencies by execing the node license-checker script
func getJavascriptDependencies() ([]Dependency, error) {
	out, err := exec.Command("node_modules/license-checker/bin/license-checker", "--csv").Output()
	if err != nil {
		return nil, errors.Wrap(err, "running license-checker")
	}

	reader := csv.NewReader(bytes.NewReader(out))

	fields, err := reader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "reading fields")
	}
	if !reflect.DeepEqual(fields, []string{"module name", "license", "repository"}) {
		return nil, errors.Wrap(err, "unexpected fields")
	}

	packages, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "reading lines")
	}

	var deps []Dependency
	for _, p := range packages {
		dep := Dependency{
			Name:       p[0],
			License:    p[1],
			Repository: p[2],
		}

		// When license-checker can't find a 'license' key in the
		// package.json, it "guesses" the license by looking for a
		// LICENSE or COPYING file. If the license was inferred by
		// that, there will be a trailing '*'
		if dep.License[len(dep.License)-1] == '*' {
			dep.License = dep.License[:len(dep.License)-1]
		}

		deps = append(deps, dep)
	}

	return deps, nil
}

func isLicenseCompatible(settings Settings, dep Dependency) bool {
	// Manual exception by name
	if _, ok := settings.Exceptions[dep.Name]; ok {
		return true
	}

	// License matches allowable licenses
	if _, ok := settings.AllowedLicenses[dep.License]; ok {
		return true
	}

	return false
}

func checkLicenses(settings Settings, deps []Dependency) []Dependency {
	var incompatible []Dependency
	for _, dep := range deps {
		if !isLicenseCompatible(settings, dep) {
			incompatible = append(incompatible, dep)
		}
	}

	sort.Slice(incompatible, func(i, j int) bool {
		return strings.ToLower(incompatible[i].Name) <= strings.ToLower(incompatible[j].Name)
	})

	sort.SliceStable(incompatible, func(i, j int) bool {
		return strings.ToLower(incompatible[i].License) <= strings.ToLower(incompatible[j].License)
	})

	return incompatible
}

func main() {
	settingsContents, err := ioutil.ReadFile("lint_license/license_settings.yaml")
	if err != nil {
		log.Fatal("error reading settings file: ", err)
	}

	var settings Settings
	err = yaml.Unmarshal(settingsContents, &settings)
	if err != nil {
		log.Fatal("error unmarshaling settings: ", err)
	}

	fmt.Println("Retrieving JS dependencies")

	jsDeps, err := getJavascriptDependencies()
	if err != nil {
		log.Fatal("error retrieving JS deps: ", err)
	}

	fmt.Printf("Checking %d JS dependencies\n", len(jsDeps))

	incompatibleJS := checkLicenses(settings, jsDeps)

	fmt.Printf("Found %d incompatible licenses\n", len(incompatibleJS))

	if len(incompatibleJS) > 0 {
		for _, dep := range incompatibleJS {
			fmt.Printf("Found incompatible license '%s' for dependency '%s'\n",
				dep.License, dep.Name)
		}
	}

	if len(incompatibleJS) > 0 {
		os.Exit(1)
	}
}
