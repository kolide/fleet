// package logging provides logger "plugins" for writing osquery status and
// result logs to various destinations.
package logging

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

type OsqueryLogger struct {
	Status kolide.JSONLogger
	Result kolide.JSONLogger
}

func New(config config.KolideConfig, logger log.Logger) (*OsqueryLogger, error) {
	var status, result kolide.JSONLogger
	var err error

	switch config.Osquery.StatusLogPlugin {
	case "":
		// Allow "" to mean filesystem for backwards compatibility
		level.Info(logger).Log("msg", "kolide_status_log_plugin not explicitly specified. Assuming 'filesystem'")
		fallthrough
	case "filesystem":
		status, err = NewFilesystemLogWriter(
			config.Filesystem.StatusLogFile,
			logger,
			config.Filesystem.EnableLogRotation,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create filesystem status logger")
		}
	case "firehose":
		status, err = NewFirehoseLogWriter(
			config.Firehose.Region,
			config.Firehose.AccessKeyID,
			config.Firehose.SecretAccessKey,
			config.Firehose.StatusStream,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create firehose status logger")
		}
	case "pubsub":
		status, err = NewPubSubLogWriter(
			config.PubSub.Project,
			config.PubSub.StatusTopic,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create pubsub status logger")
		}
	case "stdout":
		status, err = NewStdoutLogWriter()
		if err != nil {
			return nil, errors.Wrap(err, "create stdout status logger")
		}
	default:
		return nil, errors.Errorf(
			"unknown status log plugin: %s", config.Osquery.StatusLogPlugin,
		)
	}

	switch config.Osquery.ResultLogPlugin {
	case "":
		// Allow "" to mean filesystem for backwards compatibility
		level.Info(logger).Log("msg", "kolide_result_log_plugin not explicitly specified. Assuming 'filesystem'")
		fallthrough
	case "filesystem":
		result, err = NewFilesystemLogWriter(
			config.Filesystem.ResultLogFile,
			logger,
			config.Filesystem.EnableLogRotation,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create filesystem result logger")
		}
	case "firehose":
		result, err = NewFirehoseLogWriter(
			config.Firehose.Region,
			config.Firehose.AccessKeyID,
			config.Firehose.SecretAccessKey,
			config.Firehose.ResultStream,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create firehose result logger")
		}
	case "pubsub":
		result, err = NewPubSubLogWriter(
			config.PubSub.Project,
			config.PubSub.ResultTopic,
			logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create pubsub result logger")
		}
	case "stdout":
		result, err = NewStdoutLogWriter()
		if err != nil {
			return nil, errors.Wrap(err, "create stdout result logger")
		}
	default:
		return nil, errors.Errorf(
			"unknown result log plugin: %s", config.Osquery.StatusLogPlugin,
		)
	}
	return &OsqueryLogger{Status: status, Result: result}, nil
}
