package service

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (vm validationMiddleware) ImportConfig(ctx context.Context, cfg *kolide.ImportConfig) (*kolide.ImportConfigResponse, error) {
	var invalid invalidArgumentError
	vm.validateConfigOptions(cfg, &invalid)
	vm.validatePacks(cfg, &invalid)

	return vm.Service.ImportConfig(ctx, cfg)
}

func (vm validationMiddleware) validateConfigOptions(cfg *kolide.ImportConfig, argErrs *invalidArgumentError) {

	for optName, optValue := range cfg.Options {
		opt, err := vm.ds.OptionByName(string(optName))
		if err != nil {
			argErrs.Append("options", fmt.Sprintf("invalid name '%s'", optName))
			continue
		}
		if !opt.SameType(optValue) {
			argErrs.Append("options", fmt.Sprintf("invalid type for '%s'", optName))
		}

	}
}

func (vm validationMiddleware) validatePacks(cfg *kolide.ImportConfig, argErrs *invalidArgumentError) {
	for packName, pack := range cfg.Packs {
		// if glob packs is defined we expect at least one external pack
		if packName == kolide.GlobPacks {
			if len(cfg.GlobPackNames) == 0 {
				argErrs.Append("external_packs", "missing glob packs")
				continue
			}
			// make sure that each glob pack has JSON content
			for _, p := range cfg.GlobPackNames {
				if _, ok := cfg.ExternalPacks[p]; !ok {
					argErrs.Append("external_packs", fmt.Sprintf("missing content for '%s'", p))
				}
			}
			continue
		}
		// if value is a string we expect a file path, in this case, the user has to supply the
		// contents of said file which we store in ExternalPacks, if it's not there we need to
		// raise an error
		switch pack.(type) {
		case string:
			if _, ok := cfg.ExternalPacks[packName]; !ok {
				argErrs.Append("external_packs", fmt.Sprintf("missing content for '%s'", packName))
			}
		}
	}
}
