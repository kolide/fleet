package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	license "github.com/ryanuber/go-license"

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
	Version    string
	Path       string
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

type packageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Most packages store license info in 'license'
	License interface{} `json:"license"`
	// A few store license info in an array in 'licenses'
	Licenses []interface{} `json:"licenses"`
}

func extractJSPackageInfo(path string) (Dependency, error) {
	dep := Dependency{
		Path: filepath.Dir(path),
	}

	f, err := os.Open(path)
	if err != nil {
		return dep, errors.Wrap(err, "opening package.json")
	}

	var pkg packageJSON
	err = json.NewDecoder(f).Decode(&pkg)
	if err != nil {
		return dep, errors.Wrap(err, "reading JSON from package.json")
	}

	dep.Name = pkg.Name
	dep.Version = pkg.Version

	// Pick whichever top-level license key we found
	var licObj interface{}
	if pkg.License != nil {
		licObj = pkg.License
	} else if len(pkg.Licenses) > 0 {
		licObj = pkg.Licenses[0]
	}

	switch lic := licObj.(type) {
	case string:
		// Almost all use a string value for license
		dep.License = lic

	case map[string]interface{}:
		// Some few packages use a map with the key 'type'
		// corresponding to the license name
		if lic, ok := lic["type"].(string); ok {
			dep.License = lic
		}
	}

	// If finding license info in package.json fails, we can try to
	// identify the license with the go-license package
	if dep.License == "" {
		if l, err := license.NewFromDir(dep.Path); err == nil {
			dep.License = l.Type
		}
	}

	return dep, nil
}

func getJSDeps() ([]Dependency, error) {
	var deps []Dependency
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading path %s: %s\n", path, err.Error())
		}

		// Traversal can ignore anything that is not package.json
		if info.IsDir() || info.Name() != "package.json" {
			return nil
		}

		dep, err := extractJSPackageInfo(path)
		if err != nil {
			fmt.Printf("Error analyzing path %s: %s\n", path, err.Error())
		}
		deps = append(deps, dep)

		return nil
	}

	fmt.Println("starting walk")
	t := time.Now()
	err := filepath.Walk("./node_modules", walkFn)
	fmt.Println("completed in ", time.Now().Sub(t))

	if err != nil {
		return nil, errors.Wrap(err, "walking node_modules")
	}

	return deps, nil
}

func isLicenseCompatible(settings Settings, dep Dependency) bool {
	// Manual exception by path
	if _, ok := settings.Exceptions[dep.Path]; ok {
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

	jsDeps, err := getJSDeps()
	if err != nil {
		log.Fatal("error retrieving JS deps: ", err)
	}

	fmt.Printf("Checking %d JS dependencies\n", len(jsDeps))

	incompatibleJS := checkLicenses(settings, jsDeps)

	fmt.Printf("Found %d incompatible licenses\n", len(incompatibleJS))

	if len(incompatibleJS) > 0 {
		for _, dep := range incompatibleJS {
			fmt.Printf("Incompatible license '%s' for dependency '%s' (path '%s')\n",
				dep.License, dep.Name, dep.Path)
		}
	}

	if len(incompatibleJS) > 0 {
		os.Exit(1)
	}
}
