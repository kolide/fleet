package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

		// JS packages should always have a package.json
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

type glideImport struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// glideLock is a yaml schema for the relevant portions of glide.lock
type glideLock struct {
	Imports []glideImport `yaml:"imports"`
}

func extractGoPackageInfo(pkg glideImport) (Dependency, error) {
	dep := Dependency{
		Path:    filepath.Join("vendor", pkg.Name),
		Name:    pkg.Name,
		Version: pkg.Version,
	}

	if l, err := license.NewFromDir(dep.Path); err == nil {
		dep.License = l.Type
	}

	return dep, nil
}

func getGoDeps() ([]Dependency, error) {
	glockContents, err := ioutil.ReadFile("glide.lock")
	if err != nil {
		return nil, errors.Wrap(err, "error reading glide.lock")
	}

	var glock glideLock
	err = yaml.Unmarshal(glockContents, &glock)
	if err != nil {
		log.Fatal("error unmarshaling settings: ", err)
	}

	var deps []Dependency
	for _, pkg := range glock.Imports {
		dep, err := extractGoPackageInfo(pkg)
		if err == nil {
			deps = append(deps, dep)
		}
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

	// Javascript
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

	fmt.Printf("\n\n")

	// Go
	fmt.Println("Retrieving Go dependencies")

	goDeps, err := getGoDeps()
	if err != nil {
		log.Fatal("error retrieving GO deps: ", err)
	}

	fmt.Printf("Checking %d GO dependencies\n", len(goDeps))

	incompatibleGo := checkLicenses(settings, goDeps)

	fmt.Printf("Found %d incompatible licenses\n", len(incompatibleGo))

	if len(incompatibleGo) > 0 {
		for _, dep := range incompatibleGo {
			fmt.Printf("Incompatible license '%s' for dependency '%s' (path '%s')\n",
				dep.License, dep.Name, dep.Path)
		}
	}

	if len(incompatibleJS) > 0 || len(incompatibleGo) > 0 {
		os.Exit(1)
	}
}
