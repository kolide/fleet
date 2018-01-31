package logwriter

import (
	"context"
	"os"

	"github.com/kolide/fleet/server/config"
	hostctx "github.com/kolide/fleet/server/contexts/host"

	"cloud.google.com/go/logging"
	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type stackdriverWriter struct {
	client      *logging.Client
	stackdriver *logging.Logger
	logger      kitlog.Logger
}

func newStackDriver(conf config.OsqueryLogStackDriver, logger kitlog.Logger) (*stackdriverWriter, error) {
	if conf.ProjectID == "" && conf.ApplicationCredentials == "" {
		return nil, nil
	}

	// XXX(thorduri): Banana.
	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", conf.ApplicationCredentials); err != nil {
		return nil, err
	}

	client, err := logging.NewClient(context.Background(), conf.ProjectID)
	if err != nil {
		return nil, err
	}

	// XXX(thorduri): Status/Result annotations.
	return &stackdriverWriter{
		client:      client,
		stackdriver: client.Logger("osquery"),
		logger:      kitlog.With(logger, "type", "stackdriver"),
	}, nil
}

func (sdw *stackdriverWriter) Write(ctx context.Context, b []byte) error {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		// this can only happen if the host authz failed.
		sdw.stackdriver.Log(logging.Entry{
			Severity: logging.Alert,
			Payload:  "logging osquery log without host ctx",
		})
		return nil
	}

	sdw.stackdriver.Log(logging.Entry{
		Labels: map[string]string{
			"osquery_host_hostname": host.HostName,
			"osquery_host_uuid":     host.UUID,
		},
		Payload: b,
	})

	return nil
}

func (sdw *stackdriverWriter) Flush(ctx context.Context) error {
	if err := sdw.stackdriver.Flush(); err != nil {
		return errors.Wrap(err, "stackdriverwriter: flush failed")
	}
	return nil
}

func (sdw *stackdriverWriter) Close(ctx context.Context) error {
	if err := sdw.client.Close(); err != nil {
		return errors.Wrap(err, "stackdriverwriter: close failed")
	}
	return nil
}
