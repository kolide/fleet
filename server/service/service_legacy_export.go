package service

import (
	"context"
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

type configForExport struct {
	Options   map[string]interface{} `json:"options"`
	FilePaths map[string][]string    `json:"file_paths,omitempty"`
}

type yamlObjForExport struct {
	ApiVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Spec       specForExport `json:"spec"`
}

type specForExport struct {
	Config json.RawMessage `json:"config"`
}

func (svc service) ExportConfig(ctx context.Context) (string, error) {
	options, err := svc.ds.GetOsqueryConfigOptions()
	if err != nil {
		return "", errors.Wrap(err, "getting osquery options")
	}

	fimConfig, err := svc.GetFIM(ctx)
	if err != nil {
		return "", errors.Wrap(err, "getting FIM configs")
	}

	config := configForExport{
		Options:   options,
		FilePaths: fimConfig.FilePaths,
	}

	confJSON, err := json.Marshal(config)

	confObj := yamlObjForExport{
		ApiVersion: "TODO",
		Kind:       "TODO",
		Spec: specForExport{
			Config: json.RawMessage(confJSON),
		},
	}

	confYaml, err := yaml.Marshal(confObj)
	if err != nil {
		return "", errors.Wrap(err, "marshal YAML")
	}

	return string(confYaml), nil
}
