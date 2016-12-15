package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockValidationItem struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
type mockValidationError struct {
	Message string               `json:"message"`
	Errors  []mockValidationItem `json:"errors"`
}

func testGetAppConfig(t *testing.T, r *testResource) {

	req, err := http.NewRequest("GET", r.server.URL+"/api/v1/kolide/config", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	var configInfo getAppConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&configInfo)
	require.Nil(t, err)
	require.NotNil(t, configInfo.AppConfig)
	config := configInfo.AppConfig
	assert.Equal(t, uint(465), config.SMTPPort)
	assert.Equal(t, "Kolide", config.OrgName)
	assert.Equal(t, "http://foo.bar/image.png", config.OrgLogoURL)

}

func testModifyAppConfig(t *testing.T, r *testResource) {
	body := kolide.ModifyAppConfigRequest{
		TestSMTP: false,
		AppConfig: kolide.AppConfig{
			KolideServerURL:        "https://foo.com",
			OrgName:                "Zip",
			OrgLogoURL:             "http://foo.bar/image.png",
			SMTPPort:               567,
			SMTPAuthenticationType: kolide.AuthTypeNone,
			SMTPServer:             "foo.com",
			SMTPEnableTLS:          true,
			SMTPVerifySSLCerts:     true,
			SMTPEnableStartTLS:     true,
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	require.Nil(t, err)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/config", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var respBody modifyAppConfigResponse

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.Nil(t, err)
	assert.Equal(t, body.AppConfig.OrgName, respBody.Response.OrgName)

}

func testModifyAppConfigWithValidationFail(t *testing.T, r *testResource) {

	body := kolide.ModifyAppConfigRequest{
		TestSMTP: false,
		AppConfig: kolide.AppConfig{
			OrgName:                "Zip",
			OrgLogoURL:             "http://foo.bar/image.png",
			SMTPPort:               567,
			SMTPAuthenticationType: kolide.AuthTypeNone,
			SMTPEnableTLS:          true,
			SMTPVerifySSLCerts:     true,
			SMTPEnableStartTLS:     true,
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	require.Nil(t, err)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/config", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var validationErrors mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&validationErrors)
	require.Nil(t, err)
	assert.Equal(t, "Validation Failed", validationErrors.Message)
	assert.Equal(t, 2, len(validationErrors.Errors))
	assert.Equal(t, "kolide_server_url", validationErrors.Errors[0].Name)
	assert.Equal(t, "url scheme must be https", validationErrors.Errors[0].Reason)
	assert.Equal(t, "smtp_server", validationErrors.Errors[1].Name)
	assert.Equal(t, "missing require argument", validationErrors.Errors[1].Reason)
	// verify no changes are not saved if validation fails
	config, _ := r.ds.AppConfig()
	assert.NotEqual(t, config.OrgName, body.AppConfig.OrgName)
}
