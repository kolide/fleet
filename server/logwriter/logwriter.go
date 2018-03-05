// Package logwriter provides logging utilities for writing osquery results and status logs
package logwriter

import (
	"context"

	"github.com/kolide/fleet/server/config"

	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type writers interface {
	Write(context.Context, []byte) error
	Flush(context.Context) error
	Close(context.Context) error
}

type Log struct {
	writers []writers
	logger  kitlog.Logger
}

func New(conf config.OsqueryLog, logger kitlog.Logger) (*Log, error) {
	l := new(Log)

	logger = kitlog.With(logger, "component", "logwriter")

	fw, err := newFile(conf.File, logger)
	if err != nil {
		return nil, err
	} else if fw != nil {
		l.writers = append(l.writers, fw)
	}

	sw, err := newSumo(conf.Sumo, logger)
	if err != nil {
		return nil, err
	} else if sw != nil {
		l.writers = append(l.writers, sw)
	}

	sdw, err := newStackDriver(conf.StackDriver, logger)
	if err != nil {
		return nil, err
	} else if sdw != nil {
		l.writers = append(l.writers, sdw)
	}

	l.logger = logger
	return l, nil
}

func (l *Log) Write(ctx context.Context, b []byte) error {
	nerr := 0

	for _, w := range l.writers {
		if err := w.Write(ctx, b); err != nil {
			l.logger.Log("err", err)
			nerr++
		}
	}

	// Only report errors if all writers failed.
	if nerr == len(l.writers) {
		return errors.New("logwriter: all writers failed writing")
	}

	return nil
}

func (l *Log) Flush(ctx context.Context) error {
	nerr := 0

	for _, w := range l.writers {
		if err := w.Flush(ctx); err != nil {
			l.logger.Log("err", err)
			nerr++
		}
	}

	// Only report errors if all writers failed.
	if nerr == len(l.writers) {
		return errors.New("logwriter: all writers failed flushing")
	}

	return nil
}

func (l *Log) Close(ctx context.Context) error {
	nerr := 0

	for _, w := range l.writers {
		if err := w.Close(ctx); err != nil {
			l.logger.Log("err", err)
			nerr++
		}
	}

	// Only report errors if all writers failed.
	if nerr == len(l.writers) {
		return errors.New("logwriter: all writers failed flushing")
	}

	return nil
}
