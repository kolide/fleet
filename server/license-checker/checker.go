package license

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kolide/server/kolide"
)

const defaultPollFrequency = time.Hour
const defaultHttpClientTimeout = 10 * time.Second

type revokeInfo struct {
	UUID    string `json:"uuid"`
	Revoked bool   `json:"revoked"`
}

type revokeError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// Checker checks remote kolide/cloud app for license revocation
// status
type Checker struct {
	ds                kolide.Datastore
	logger            log.Logger
	url               string
	pollFrequency     time.Duration
	httpClientTimeout time.Duration
	finish            chan struct{}
}

type Option func(chk *Checker)

// Logger set the logger that will be used by the Checker
func Logger(logger log.Logger) Option {
	return func(chk *Checker) {
		chk.logger = logger
	}
}

// PollFrequency defines frequency to check for license revocation.
// Default is once per hour.
func PollFrequency(freq time.Duration) Option {
	return func(chk *Checker) {
		chk.pollFrequency = freq
	}
}

// HTTPClientTimeout determines how long to wait for requests to remote
// host to response.  Defaults to 10 seconds
func HTTPClientTimeout(timeout time.Duration) Option {
	return func(chk *Checker) {
		chk.httpClientTimeout = timeout
	}
}

// NewChecker instantiates a service that will check periodically to see if a license
// is revoked.  licenseEndpointURL is the root url for kolide/cloud server.  For example
// https://cloud.kolide.co/api/v0/licenses
// You may optionally set a logger, and/or supply a polling frequency that defines
// how often we check for revocation.
func NewChecker(ds kolide.Datastore, licenseEndpointURL string, opts ...Option) *Checker {
	response := &Checker{
		pollFrequency:     defaultPollFrequency,
		httpClientTimeout: defaultHttpClientTimeout,
		logger:            log.NewNopLogger(),
		ds:                ds,
		url:               licenseEndpointURL,
	}
	for _, o := range opts {
		o(response)
	}
	response.logger = log.NewContext(response.logger).With("component", "license-checker")
	return response
}

// Start begins checking for license revocation
func (cc *Checker) Start() {
	cc.finish = make(chan struct{})
	// pass in copy of receiver to avoid race conditions
	go func(chk Checker) {
		chk.logger.Log("msg", "starting")
		for {
			select {
			case <-chk.finish:
				chk.logger.Log("msg", "finishing")
				return
			case <-time.After(chk.pollFrequency):
				updateLicenseRevocation(&chk)
			}
		}
	}(*cc)
}

// Stop ends checking for license revocation.
func (cc *Checker) Stop() {
	if cc.finish != nil {
		close(cc.finish)
		cc.finish = nil
	}
}

func updateLicenseRevocation(chk *Checker) {
	chk.logger.Log("msg", "begin license check")
	defer chk.logger.Log("msg", "ending license check")

	license, err := chk.ds.License()
	if err != nil {
		chk.logger.Log("msg", "couldn't fetch license", "err", err)
		return
	}
	claims, err := license.Claims()
	if err != nil {
		chk.logger.Log("msg", "fetching claims", "err", err)
		return
	}
	url := fmt.Sprintf("%s/%s", chk.url, claims.LicenseUUID)
	client := http.Client{Timeout: chk.httpClientTimeout}
	resp, err := client.Get(url)
	if err != nil {
		chk.logger.Log("msg", fmt.Sprintf("fetching %s", url), "err", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var revInfo revokeInfo
		err = json.NewDecoder(resp.Body).Decode(&revInfo)
		if err != nil {
			chk.logger.Log("msg", "decoding response", "err", err)
			return
		}
		err = chk.ds.RevokeLicense(revInfo.Revoked)
		if err != nil {
			chk.logger.Log("msg", "revoke status", "err", err)
			return
		}
		// success
		chk.logger.Log("msg", fmt.Sprintf("license revocation status retrieved succesfully, revoked: %t", revInfo.Revoked))
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		var revInfo revokeError
		err = json.NewDecoder(resp.Body).Decode(&revInfo)
		if err != nil {
			chk.logger.Log("msg", "decoding response", "err", err)
			return
		}
		chk.logger.Log("msg", "host response", "err", fmt.Sprintf("status: %d error: %s", revInfo.Status, revInfo.Error))
		return
	}
	chk.logger.Log("msg", "host response", "err", fmt.Sprintf("unexpected response status from host, status %s", resp.Status))
}
