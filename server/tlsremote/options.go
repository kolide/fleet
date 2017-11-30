package tlsremote

import (
	"io"
	"time"

	"github.com/go-kit/kit/log"
)

// Option is a config parameter of the OsqueryService.
type Option func(*OsqueryService)

// WithLogger configures the OsqueryService with a custom logger.
func WithLogger(logger log.Logger) Option {
	return func(svc *OsqueryService) {
		svc.logger = logger
	}
}

func WithStatusLogWriters(writers ...io.Writer) Option {
	return func(svc *OsqueryService) {
		svc.osqueryStatusLogWriter = io.MultiWriter(writers...)
	}
}

func WithResultLogWriters(writers ...io.Writer) Option {
	return func(svc *OsqueryService) {
		svc.osqueryResultLogWriter = io.MultiWriter(writers...)
	}
}

func WithLableUpdateInterval(interval time.Duration) Option {
	return func(svc *OsqueryService) {
		svc.labelUpdateInterval = interval
	}
}

func WithNodeKeySize(size int) Option {
	return func(svc *OsqueryService) {
		svc.nodeKeySize = size
	}
}
