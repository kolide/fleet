package service

import (
	"github.com/kolide/kolide/server/kolide"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (svc service) GetOptions(ctx context.Context) ([]kolide.Option, error) {
	opts, err := svc.ds.ListOptions()
	if err != nil {
		return nil, errors.Wrap(err, "options service")
	}
	return opts, nil
}

func (svc service) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	if err := svc.ds.SaveOptions(req.Options); err != nil {
		return nil, errors.Wrap(err, "modify options service")
	}
	return req.Options, nil
}

func (svc service) ExpectedCheckinInterval(ctx context.Context) (uint, error) {
	interval := uint(0)
	found := false

	for _, option := range kolide.OsqueryIntervalOptionNames {
		// for each option which is known to hold a TLS check-in interval, try to
		// fetch it
		opt, err := svc.ds.OptionByName(option)
		if err != nil {
			// if the option is not set, try the next known option
			if _, ok := err.(kolide.NotFoundError); ok {
				continue
			}
			// if some other error occured when getting the option, we want to return
			// that
			return 0, err
		}

		// try to cast the option as a uint. if this fails, the option has likely been set incorrectly
		val, ok := opt.Value.Val.(uint)
		if !ok {
			return 0, errors.New("Option is not a uint: " + opt.Name)
		}

		// If an option has not been found yet, we want to save this interval.
		// If an option HAS been found already and this one is less, we want to
		// save that as our new minimum check-in interval.
		if !found || val < interval {
			found = true
			interval = val
		}
	}

	// if we never found any interval options set, the default distributed
	// interval is 60, so we use that
	if !found {
		return 60, nil
	}

	// return the lowest interval that we found
	return interval, nil
}
