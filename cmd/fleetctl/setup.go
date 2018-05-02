package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/service"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func setupCommand() cli.Command {
	var (
		flAddress            string
		flEmail              string
		flPassword           string
		flOrgName            string
		flInsecureSkipVerify bool
		flDebug              bool
	)
	return cli.Command{
		Name:      "setup",
		Usage:     "Setup a Kolide Fleet instance",
		UsageText: `fleetctl config login [options]`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "address",
				EnvVar:      "ADDRESS",
				Value:       "",
				Destination: &flAddress,
				Usage:       "The address of the Kolide Fleet instance",
			},
			cli.StringFlag{
				Name:        "email",
				EnvVar:      "EMAIL",
				Value:       "",
				Destination: &flEmail,
				Usage:       "The email of the admin user to create",
			},
			cli.StringFlag{
				Name:        "password",
				EnvVar:      "PASSWORD",
				Value:       "",
				Destination: &flPassword,
				Usage:       "The password for the admin user",
			},
			cli.StringFlag{
				Name:        "org-name",
				EnvVar:      "ORG_NAME",
				Value:       "",
				Destination: &flOrgName,
				Usage:       "The name of the organization",
			},
			cli.BoolFlag{
				Name:        "insecure-skip-verify",
				EnvVar:      "INSECURE_SKIP_VERIFY",
				Destination: &flInsecureSkipVerify,
				Usage:       "Whether or not to validate the remote TLS certificate",
			},
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			if flAddress == "" {
				return errors.New("--address is required")
			}

			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: flInsecureSkipVerify},
				},
			}

			t := true
			body := service.SetupRequest{
				Admin: &kolide.UserPayload{
					Admin:    &t,
					Username: &flEmail,
					Email:    &flEmail,
					Password: &flPassword,
				},
				OrgInfo: &kolide.OrgInfo{
					OrgName: &flOrgName,
				},
				KolideServerURL: &flAddress,
			}

			b, err := json.Marshal(body)
			if err != nil {
				return errors.Wrap(err, "error marshaling json")
			}

			request, err := http.NewRequest(
				"POST",
				flAddress+"/api/v1/setup",
				bytes.NewBuffer(b),
			)
			if err != nil {
				return errors.Wrap(err, "error creating request object")
			}
			request.Header.Set("content-type", "application/json")
			request.Header.Set("accept", "application/json")

			response, err := httpClient.Do(request)
			if err != nil {
				return errors.Wrap(err, "error making request")
			}
			defer response.Body.Close()

			// If setup has already been completed, Kolide Fleet will not serve the
			// setup route, which will cause the request to 404
			if response.StatusCode == http.StatusNotFound {
				return errors.New("Kolide Fleet has already been setup")
			}

			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("Received HTTP %d instead of HTTP 200", response.StatusCode)
			}

			responeBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return errors.Wrap(err, "error reading response body")
			}

			var responseBody service.SetupResponse
			err = json.Unmarshal(responeBytes, &responseBody)
			if err != nil {
				return errors.Wrap(err, "error decoding HTTP response body")
			}

			if responseBody.Err != nil {
				return errors.Wrap(err, "error setting up fleet instance")
			}

			// TODO save the token in ~/.fleet/config
			fmt.Println(*responseBody.Token)

			return nil
		},
	}
}
