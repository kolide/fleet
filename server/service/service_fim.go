package service

import (
	"context"

	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
	"encoding/json"
)

func (svc service) GetFIM(ctx context.Context) (*kolide.FIMConfig, error) {
	config, err := svc.ds.AppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting fim config")
	}
	paths, err := svc.ds.FIMSections()
	if err != nil {
		return nil, errors.Wrap(err, "getting fim paths")
	}

	var arr []string

	_ = json.Unmarshal([]byte(config.FIMFileAccesses), &arr)


	result := &kolide.FIMConfig{
		Interval:  uint(config.FIMInterval),
		FilePaths: paths,
		FileAccesses: arr,
	}
	return result, nil
}

// ModifyFIM will remove existing FIM settings and replace it
func (svc service) ModifyFIM(ctx context.Context, fim kolide.FIMConfig) error {
	if err := svc.ds.ClearFIMSections(); err != nil {
		return errors.Wrap(err, "updating fim")
	}
	config, err := svc.ds.AppConfig()
	if err != nil {
		return errors.Wrap(err, "updating fim")
	}

	config.FIMInterval = int(fim.Interval)
	fileAccesses, _ := json.Marshal(fim.FileAccesses)
	config.FIMFileAccesses = string(fileAccesses)

	for sectionName, paths := range fim.FilePaths {
		section := kolide.FIMSection{
			SectionName: sectionName,
			Paths:       paths,
		}
		if _, err := svc.ds.NewFIMSection(&section); err != nil {
			return errors.Wrap(err, "creating fim section")
		}
	}
	return svc.ds.SaveAppConfig(config)
}
