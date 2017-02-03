package license

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kolide/server/kolide"
)

const defaultPollFrequency = time.Hour

type Checker interface {
	// Start begins checking for license revocation
	Start()
	// Stop ends checking
	Stop()
}

type checker struct {
	ds            kolide.Datastore
	logger        log.Logger
	url           string
	pollFrequency time.Duration
	finish        chan struct{}
}

type Option func(opts *options)

type options struct {
	logger        log.Logger
	pollFrequency time.Duration
}

// Logger set the logger that will be used by the Checker
func Logger(logger log.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

// PollFrequency defines frequency to check for license revocation
func PollFrequency(freq time.Duration) Option {
	return func(opt *options) {
		opt.pollFrequency = freq
	}
}

// NewChecker instantiates a service that will check periodically to see if a license
// is revoked.  licenseEndpointURL is the root url for kolide/cloud server.  For example
// https://cloud.kolide.co/api/v0/licenses
// You may optionally set a logger, and/or supply a polling frequency that defines
// how often we check for revocation.
func NewChecker(ds kolide.Datastore, licenseEndpointURL string, opts ...Option) Checker {
	settings := options{pollFrequency: defaultPollFrequency}
	for _, o := range opts {
		o(&settings)
	}
	if settings.logger == nil {
		w := log.NewSyncWriter(os.Stderr)
		settings.logger = log.NewLogfmtLogger(w)
	}

	return &checker{
		ds:            ds,
		logger:        settings.logger,
		pollFrequency: settings.pollFrequency,
		url:           licenseEndpointURL,
	}
}

func (cc *checker) Start() {
	cc.finish = make(chan struct{})
	go func(chk checker) {
		logMsg(cc, "starting")
		for {

			select {
			case <-chk.finish:
				logMsg(cc, "finishing")
				return
			case <-time.After(chk.pollFrequency):
				return
			}
		}
	}(*cc)
}

func (cc *checker) Stop() {
	if cc.finish != nil {
		close(cc.finish)
		cc.finish = nil
	}
}

func updateLicenseRevocation(chk *checker) {
	logMsg(chk, "begin license check")
	defer logMsg(chk, "ending license check")

	license, err := chk.ds.License()
	if err != nil {
		logErr(chk, "couldn't fetch license", err.Error())
	}
	if license.Token == nil {
		logMsg(chk, "no license present")
		return
	}

}

func logMsg(chk *checker, msg string) {
	chk.logger.Log("component", "license-checker", "msg", msg)
}

func logErr(chk *checker, msg, errMsg string) {
	chk.logger.Log("component", "license-checker", "msg", msg, "err", errMsg)
}
