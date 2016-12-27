package service

import (
	"fmt"
	"strconv"

	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (mw validationMiddleware) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	transformed := kolide.OptionRequest{Options: []kolide.Option{}}
	invalid := &invalidArgumentError{}
	for _, opt := range req.Options {
		if opt.ReadOnly {
			invalid.Append(opt.Name, "readonly option")
			continue
		}
		val, err := optValToDB(opt.Value, opt.Type)
		if err != nil {
			invalid.Append(opt.Name, err.Error())
			continue
		}
		opt.Value = val
		transformed.Options = append(transformed.Options, opt)
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyOptions(ctx, transformed)
}

var errTypeMismatch = fmt.Errorf("type mismatch")
var errInvalidType = fmt.Errorf("invalid option type")

func optValToDB(v interface{}, typ kolide.OptionType) (interface{}, error) {
	var str string
	if v == nil {
		return v, nil
	}

	switch typ {
	case kolide.OptionTypeFlag:
		flag, ok := v.(bool)
		if !ok {
			return nil, errTypeMismatch
		}
		str = strconv.FormatBool(flag)
	case kolide.OptionTypeString:
		s, ok := v.(string)
		if !ok {
			return nil, errTypeMismatch
		}
		str = s
	case kolide.OptionTypeInt:
		num, ok := v.(float64)
		if !ok {
			return nil, errTypeMismatch
		}
		str = strconv.Itoa(int(num))
	default:
		return nil, errInvalidType
	}
	return &str, nil
}
