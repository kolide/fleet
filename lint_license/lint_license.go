package lint_license

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os/exec"
	"reflect"
)

type pkgInfo []string

func (p pkgInfo) Name() string {
	return p[0]
}

func (p pkgInfo) License() string {
	return p[1]
}

func (p pkgInfo) Repository() string {
	return p[2]
}

var allowedLicenses = map[string]bool{
	"MIT":        true,
	"BSD":        true,
	"Apache-1.0": true,
	"JSON":       true,
	"Postgres":   true,
	"ISC":        true,
	"Apache-2.0": true,
}

func main() {
	out, err := exec.Command("node_modules/license-checker/bin/license-checker", "--csv").Output()
	if err != nil {
		log.Fatal("error running license-checker", err)
	}

	reader := csv.NewReader(bytes.NewReader(out))

	fields, err := reader.Read()
	if err != nil {
		log.Fatal("error reading fields: ", err)
	}
	if !reflect.DeepEqual(fields, []string{"module name", "license", "repository"}) {
		log.Fatal("unexpected fields: ", fields)
	}

	packages, err := reader.ReadAll()
	if err != nil {
		log.Fatal("error reading lines: ", err)
	}

	for _, p := range packages {
		pkg := pkgInfo(p)

		if !allowedLicenses[pkg.License()] {
			fmt.Println(pkg.Name(), pkg.License(), pkg.Repository())
		}
	}
}
