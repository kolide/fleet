// Package service holds the implementation of the kolide service interface and the HTTP endpoints
// for the API
package service

import (
	"net/http"
	"time"

	"github.com/WatchBeam/clock"
	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/logwriter"
	"github.com/kolide/fleet/server/sso"
)

// NewService creates a new service from the config struct
func NewService(ds kolide.Datastore, resultStore kolide.QueryResultStore,
	logger kitlog.Logger, kolideConfig config.KolideConfig, mailService kolide.MailService,
	c clock.Clock, sso sso.SessionStore) (kolide.Service, error) {
	var svc kolide.Service

	statusLog, err := logwriter.New(kolideConfig.Osquery.StatusLog, logger)
	if err != nil {
		return nil, err
	}
	resultLog, err := logwriter.New(kolideConfig.Osquery.ResultLog, logger)
	if err != nil {
		return nil, err
	}

	svc = service{
		ds:          ds,
		resultStore: resultStore,
		logger:      logger,
		config:      kolideConfig,
		clock:       c,

		osqueryStatusLog: statusLog,
		osqueryResultLog: resultLog,
		mailService:      mailService,
		ssoSessionStore:  sso,
		metaDataClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	svc = validationMiddleware{svc, ds, sso}
	return svc, nil
}

type service struct {
	ds          kolide.Datastore
	resultStore kolide.QueryResultStore
	logger      kitlog.Logger
	config      config.KolideConfig
	clock       clock.Clock

	osqueryStatusLog *logwriter.Log
	osqueryResultLog *logwriter.Log

	mailService     kolide.MailService
	ssoSessionStore sso.SessionStore
	metaDataClient  *http.Client
}

func (s service) SendEmail(mail kolide.Email) error {
	return s.mailService.SendEmail(mail)
}

func (s service) Clock() clock.Clock {
	return s.clock
}

type validationMiddleware struct {
	kolide.Service
	ds              kolide.Datastore
	ssoSessionStore sso.SessionStore
}
