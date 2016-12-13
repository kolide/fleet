package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kolide/kolide-ose/server/kolide"
)

type mockValidationItem struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
type mockValidationError struct {
	Message string               `json:"message"`
	Errors  []mockValidationItem `json:"errors"`
}

func (s *EndpointTestSuite) TestGetAppConfig() {
	req, err := http.NewRequest("GET", s.server.URL+"/api/v1/kolide/config", nil)
	s.Require().Nil(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	s.Require().Nil(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	var configInfo getAppConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&configInfo)
	s.Require().Nil(err)
	s.Require().NotNil(configInfo.AppConfig)
	config := configInfo.AppConfig
	s.Equal(uint(465), config.Port)
	s.Equal("Kolide", config.OrgName)
	s.Equal("http://foo.bar/image.png", config.OrgLogoURL)

}

func (s *EndpointTestSuite) TestModifyAppConfig() {
	body := kolide.ModifyAppConfigRequest{
		TestSMTP: false,
		AppConfig: kolide.AppConfig{
			KolideServerURL: "https://foo.com",
			OrgName:         "Zip",
			OrgLogoURL:      "http://foo.bar/image.png",
			SMTPConfig: &kolide.SMTPConfig{
				Port:               567,
				AuthenticationType: kolide.AuthTypeNone,
				Server:             "foo.com",
				EnableTLS:          true,
				VerifySSLCerts:     true,
				EnableStartTLS:     true,
			},
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	s.Require().Nil(err)
	req, err := http.NewRequest("PATCH", s.server.URL+"/api/v1/kolide/config", &buffer)
	s.Require().Nil(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	s.Require().Nil(err)

	var respBody modifyAppConfigResponse

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	s.Require().Nil(err)
	s.Equal(body.AppConfig.OrgName, respBody.Response.OrgName)

}

func (s *EndpointTestSuite) TestModifyAppConfigWithValidationFail() {

	body := kolide.ModifyAppConfigRequest{
		TestSMTP: false,
		AppConfig: kolide.AppConfig{
			OrgName:    "Zip",
			OrgLogoURL: "http://foo.bar/image.png",
			SMTPConfig: &kolide.SMTPConfig{
				Port:               567,
				AuthenticationType: kolide.AuthTypeNone,
				EnableTLS:          true,
				VerifySSLCerts:     true,
				EnableStartTLS:     true,
			},
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	s.Require().Nil(err)
	req, err := http.NewRequest("PATCH", s.server.URL+"/api/v1/kolide/config", &buffer)
	s.Require().Nil(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	s.Require().Nil(err)

	var validationErrors mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&validationErrors)
	s.Require().Nil(err)
	s.Equal("Validation Failed", validationErrors.Message)
	s.Equal(2, len(validationErrors.Errors))
	s.Equal("kolide_server_url", validationErrors.Errors[0].Name)
	s.Equal("url scheme must be https", validationErrors.Errors[0].Reason)
	s.Equal("smtp_server", validationErrors.Errors[1].Name)
	s.Equal("missing require argument", validationErrors.Errors[1].Reason)
	// verify no changes are not saved if validation fails
	config, _ := s.ds.AppConfig()
	s.NotEqual(config.OrgName, body.AppConfig.OrgName)
}
