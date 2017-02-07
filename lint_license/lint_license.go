package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	license "github.com/ryanuber/go-license"

	yaml "gopkg.in/yaml.v2"
)

// This script is intended to be run from the root of the Kolide repo. All
// paths are relative to that directory.
const configPath = "./lint_license/license_settings.yaml"
const templatePath = "./lint_license/dependencies.md.tmpl"
const templateName = "dependencies.md.tmpl"
const glideLockPath = "./glide.lock"
const nodeModulesPath = "./node_modules"
const vendorPath = "./vendor"
const jsSourceURLBase = "https://www.npmjs.com/package/"
const generatedMarkdownPath = "./docs/third-party/dependencies.md"

// settings defines the config options for this script
type settings struct {
	// AllowedLicenses is a map from acceptable license name to the URL for
	// that license.
	AllowedLicenses map[string]string `yaml:"allowed_licenses"`
	// Overrides is a map of package paths to override license names. These
	// licenses are determined by a human and manually overridden.
	Overrides map[string]string `yaml:"overrides"`
	// Tests is a set of packages that are tests for another package, and
	// should not be counted as a separate dependency. They are determined
	// by a human and manually overridden.
	Tests map[string]struct{} `yaml:"tests"`
}

// dependency stores all the relevant info for a Kolide dependency
type dependency struct {
	// Name is the package name
	Name string
	// License is the name of the license used
	License string
	// SourceURL is the URL for the package
	SourceURL string
	// Version is the version we are using
	Version string
	// Path is the local directory path
	Path string
}

// packageJSON is a schema for the relevant bits of package.json
type packageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Most packages store license info in 'license'
	License interface{} `json:"license"`
	// A few store license info in an array in 'licenses'
	Licenses []interface{} `json:"licenses"`
}

func extractJSPackageInfo(config settings, path string) (dependency, error) {
	dep := dependency{
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
	dep.SourceURL = jsSourceURLBase + dep.Name

	if lic, ok := config.Overrides[dep.Path]; ok {
		dep.License = lic
		return dep, nil
	}

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

func getJSDeps(config settings) ([]dependency, error) {
	var deps []dependency
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading path %s: %s\n", path, err.Error())
		}

		// JS packages should always have a package.json
		if info.IsDir() || info.Name() != "package.json" {
			return nil
		}

		// Skip test packages that are explicitly excluded by the
		// config file
		if _, ok := config.Tests[filepath.Dir(path)]; ok {
			return nil
		}

		dep, err := extractJSPackageInfo(config, path)
		if err != nil {
			fmt.Printf("Error analyzing path %s: %s\n", path, err.Error())
		}
		deps = append(deps, dep)

		return nil
	}

	err := filepath.Walk(nodeModulesPath, walkFn)
	if err != nil {
		return nil, errors.Wrap(err, "walking node_modules")
	}

	return deps, nil
}

//
type glideImport struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// glideLock is a yaml schema for the relevant portions of glide.lock
type glideLock struct {
	Imports []glideImport `yaml:"imports"`
}

func extractGoPackageInfo(config settings, pkg glideImport) (dependency, error) {
	dep := dependency{
		Path:      filepath.Join(vendorPath, pkg.Name),
		Name:      pkg.Name,
		SourceURL: "https://" + pkg.Name,
		Version:   pkg.Version,
	}

	if lic, ok := config.Overrides[dep.Path]; ok {
		dep.License = lic
		return dep, nil
	}

	if l, err := license.NewFromDir(dep.Path); err == nil {
		dep.License = l.Type
	}

	return dep, nil
}

func getGoDeps(config settings) ([]dependency, error) {
	glockContents, err := ioutil.ReadFile(glideLockPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading glide.lock")
	}

	var glock glideLock
	err = yaml.Unmarshal(glockContents, &glock)
	if err != nil {
		log.Fatal("error unmarshaling glide.lock: ", err)
	}

	var deps []dependency
	for _, pkg := range glock.Imports {
		dep, err := extractGoPackageInfo(config, pkg)
		if err == nil {
			deps = append(deps, dep)
		}
	}

	return deps, nil
}

func isLicenseCompatible(config settings, dep dependency) bool {
	if _, ok := config.AllowedLicenses[dep.License]; ok {
		return true
	}

	return false
}

func checkLicenses(config settings, deps []dependency) []dependency {
	var incompatible []dependency
	for _, dep := range deps {
		if !isLicenseCompatible(config, dep) {
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

func writeDependenciesMarkdown(config settings, deps map[string]dependency, out io.Writer) error {
	funcs := template.FuncMap{
		"getLicenseURL": func(license string) string {
			return config.AllowedLicenses[license]
		},
	}

	tmpl, err := template.New("").
		Funcs(funcs).
		ParseFiles(templatePath)
	if err != nil {
		return errors.Wrap(err, "reading markdown template")
	}

	err = tmpl.ExecuteTemplate(out, templateName, deps)
	if err != nil {
		return errors.Wrap(err, "executing markdown template")
	}

	return nil
}

func main() {
	fmt.Printf("Validating dependency licenses\n\n")

	configContents, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("error reading config file: ", err)
	}

	var config settings
	err = yaml.Unmarshal(configContents, &config)
	if err != nil {
		log.Fatal("error unmarshaling config: ", err)
	}

	// Check JS deps
	fmt.Println("Retrieving JS dependencies")

	jsDeps, err := getJSDeps(config)
	if err != nil {
		log.Fatal("error retrieving JS deps: ", err)
	}

	fmt.Printf("Checking %d JS dependencies\n", len(jsDeps))

	incompatibleJS := checkLicenses(config, jsDeps)

	fmt.Printf("Found %d incompatible licenses\n", len(incompatibleJS))

	if len(incompatibleJS) > 0 {
		for _, dep := range incompatibleJS {
			fmt.Printf("Incompatible license '%s' for dependency '%s' (path '%s')\n",
				dep.License, dep.Name, dep.Path)
		}
	}

	fmt.Printf("\n")

	// Check Go deps
	fmt.Println("Retrieving Go dependencies")

	goDeps, err := getGoDeps(config)
	if err != nil {
		log.Fatal("error retrieving Go deps: ", err)
	}

	fmt.Printf("Checking %d Go dependencies\n", len(goDeps))

	incompatibleGo := checkLicenses(config, goDeps)

	fmt.Printf("Found %d incompatible licenses\n", len(incompatibleGo))

	if len(incompatibleGo) > 0 {
		for _, dep := range incompatibleGo {
			fmt.Printf("Incompatible license '%s' for dependency '%s' (path '%s')\n",
				dep.License, dep.Name, dep.Path)
		}
	}

	// Exit nonzero if incompatible licenses found
	if len(incompatibleJS) > 0 || len(incompatibleGo) > 0 {
		os.Exit(1)
	}

	// Write markdown documentation file with package/license info
	allDeps := map[string]dependency{}
	for _, dep := range jsDeps {
		allDeps[dep.Name] = dep
	}
	for _, dep := range goDeps {
		allDeps[dep.Name] = dep
	}

	out, err := os.Create(generatedMarkdownPath)
	if err != nil {
		log.Fatal("opening markdown file for writing: ", err)
	}

	fmt.Println("Writing ", generatedMarkdownPath)
	err = writeDependenciesMarkdown(config, allDeps, out)
	if err != nil {
		log.Fatal("error writing dependencies markdown: ", err)
	}
}
