package service

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kolide/fleet/server/kolide"
)

// logging middleware logs the service actions
type loggingMiddleware struct {
	kolide.Service
	logger kitlog.Logger
}

// NewLoggingService takes an existing service and adds a logging wrapper
func NewLoggingService(svc kolide.Service, logger kitlog.Logger) kolide.Service {
	return loggingMiddleware{Service: svc, logger: logger}
}

// loggerDebug returns the debug level or error if error is not nil
func (mw loggingMiddleware) loggerDebug(err error) kitlog.Logger {
	if err != nil {
		return level.Error(mw.logger)
	}
	return level.Debug(mw.logger)
}

// loggerInfo returns the info level or error if error is not nil
func (mw loggingMiddleware) loggerInfo(err error) kitlog.Logger {
	if err != nil {
		return level.Error(mw.logger)
	}
	return level.Info(mw.logger)
}
