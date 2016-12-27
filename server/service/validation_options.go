package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (mw validationMiddleware) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	invalid := &invalidArgumentError{}
	for _, opt := range req.Options {
		if opt.ReadOnly {
			invalid.Append(opt.Name, "readonly option")
		}
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyOptions(ctx, req)
}
