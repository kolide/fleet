package license

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kolide/server/kolide"
)

const defaultPollFrequency = time.Hour

// Checker checks remote kolide/cloud app for license revocation
// status
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

type revokeInfo struct {
	UUID    string `json:"uuid"`
	Revoked bool   `json:"revoked"`
}

type revokeError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
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
		logMsg(&chk, "starting")
		for {
			updateLicenseRevocation(&chk)
			select {
			case <-chk.finish:
				logMsg(&chk, "finishing")
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
		return
	}
	claims, err := license.Claims()
	if err != nil {
		logErr(chk, "fetching claims", err.Error())
		return
	}
	url := fmt.Sprintf("%s/%s", chk.url, claims.LicenseUUID)
	resp, err := http.Get(url)
	if err != nil {
		logErr(chk, fmt.Sprintf("fetching %s", url), err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var revInfo revokeInfo
		err = json.NewDecoder(resp.Body).Decode(&revInfo)
		if err != nil {
			logErr(chk, "decoding response", err.Error())
			return
		}
		err = chk.ds.RevokeLicense(revInfo.Revoked)
		if err != nil {
			logErr(chk, "revoke status", err.Error())
			return
		}
		// success
		logMsg(chk, fmt.Sprintf("license revocation status retrieved succesfully, revoked: %t", revInfo.Revoked))
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		var revInfo revokeError
		err = json.NewDecoder(resp.Body).Decode(&revInfo)
		if err != nil {
			logErr(chk, "decoding response", err.Error())
			return
		}
		logErr(chk, "host response", fmt.Sprintf("status: %d error: %s", revInfo.Status, revInfo.Error))
		return
	}
	logErr(chk, "host response", fmt.Sprintf("unexpected response status from host, status %s", resp.Status))
}

func logMsg(chk *checker, msg string) {
	chk.logger.Log("component", "license-checker", "msg", msg)
}

func logErr(chk *checker, msg, errMsg string) {
	chk.logger.Log("component", "license-checker", "msg", msg, "err", errMsg)
}
