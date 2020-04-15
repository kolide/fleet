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

// loggerForError returns a logger with a log level dependant on error
func (lm *loggingMiddleware) loggerForError(err error) kitlog.Logger {
	if err != nil {
		return level.Error(lm.logger)
	}
	return level.Debug(lm.logger)
}
