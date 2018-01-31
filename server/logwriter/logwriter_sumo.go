package logwriter

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/kolide/fleet/server/config"

	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type sumoWriter struct {
	collector string
	client    *http.Client

	buf    bytes.Buffer
	logger kitlog.Logger
	mtx    sync.Mutex
}

func newSumo(conf config.OsqueryLogSumo, logger kitlog.Logger) (*sumoWriter, error) {
	if conf.Collector == "" {
		return nil, nil
	}

	return &sumoWriter{
		collector: conf.Collector,
		logger:    kitlog.With(logger, "type", "sumo"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (sw *sumoWriter) Write(ctx context.Context, b []byte) error {
	sw.mtx.Lock()
	defer sw.mtx.Unlock()

	if _, err := sw.buf.Write(b); err != nil {
		return errors.Wrapf(err, "sumowriter: write failed")
	}
	return nil
}

func (sw *sumoWriter) Flush(ctx context.Context) error {
	var err error

	sw.mtx.Lock()
	defer sw.mtx.Unlock()

	defer func(begin time.Time) {
		_ = sw.logger.Log(
			"method", "Flush",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	// TODO(thorduri): Sumo Limitations.

	body := bytes.NewReader(sw.buf.Bytes())
	resp, err := sw.client.Post(sw.collector, "application/json", body)
	if err != nil {
		return errors.Wrap(err, "sumowriter: post failed")
	} else if resp.StatusCode != http.StatusOK {
		return errors.Errorf("sumowriter: post failed: %s", resp.Status)
	}

	sw.buf.Reset()

	return nil
}

func (sw *sumoWriter) Close(ctx context.Context) error {
	if err := sw.Flush(ctx); err != nil {
		return errors.Wrap(err, "sumowriter: close failed")
	}

	return nil
}
