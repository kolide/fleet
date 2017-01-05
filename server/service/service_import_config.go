package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) ImportConfig(ctx context.Context, cfg *kolide.ImportConfig) (*kolide.ImportConfigResponse, error) {
	resp := kolide.NewImportConfigResponse()
	if err := svc.importOptions(cfg.Options, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (svc service) importOptions(opts kolide.OptionNameToValueMap, resp *kolide.ImportConfigResponse) error {
	var updateOptions []kolide.Option
	for optName, optValue := range opts {
		opt, err := svc.ds.OptionByName(optName)
		if err != nil {
			resp.Status(kolide.OptionsSection).Warning(kolide.OptionUnknown, "skipped '%s' can't find option", optName)
			resp.Status(kolide.OptionsSection).SkipCount++
			continue
		}
		if opt.ReadOnly {
			resp.Status(kolide.OptionsSection).Warning(kolide.OptionReadonly, "skipped '%s' can't change read only option", optName)
			resp.Status(kolide.OptionsSection).SkipCount++
			continue
		}
		if opt.OptionSet() {
			resp.Status(kolide.OptionsSection).Warning(kolide.OptionAlreadySet, "skipped '%s' can't change option that is already set", optName)
			resp.Status(kolide.OptionsSection).SkipCount++
			continue
		}
		opt.SetValue(optValue)
		resp.Status(kolide.OptionsSection).Message("set %s value to %v", optName, optValue)
		resp.Status(kolide.OptionsSection).ImportCount++
		updateOptions = append(updateOptions, *opt)
	}
	if len(updateOptions) > 0 {
		if err := svc.ds.SaveOptions(updateOptions); err != nil {
			return err
		}
	}

	return nil
}
