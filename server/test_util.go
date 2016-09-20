package server

import (
	"github.com/WatchBeam/clock"
	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/kolide-ose/config"
	"github.com/kolide/kolide-ose/kolide"
)

func NewTestService(ds kolide.Datastore) (kolide.Service, error) {
	return NewService(ds, kitlog.NewNopLogger(), config.TestConfig(), nil, clock.C)
}
