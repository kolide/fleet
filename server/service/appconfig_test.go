package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func TestCleanupURL(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		name     string
	}{
		{"  http://foo.bar.com  ", "http://foo.bar.com", "leading and trailing whitespace"},
		{"\n http://foo.com \t", "http://foo.com", "whitespace"},
		{"http://foo.com", "http://foo.com", "noop"},
		{"http://foo.com/", "http://foo.com", "trailing slash"},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			actual := cleanupURL(test.in)
			assert.Equal(tt, test.expected, actual)
		})
	}

}

func TestCreateAppConfig(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)
	require.Nil(t, ds.MigrateData())

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)
	var appConfigTests = []struct {
		configPayload kolide.AppConfigPayload
	}{
		{
			configPayload: kolide.AppConfigPayload{
				OrgInfo: &kolide.OrgInfo{
					OrgLogoURL: stringPtr("acme.co/images/logo.png"),
					OrgName:    stringPtr("Acme"),
				},
				ServerSettings: &kolide.ServerSettings{
					KolideServerURL: stringPtr("https://acme.co:8080/"),
				},
			},
		},
	}

	for _, tt := range appConfigTests {
		result, err := svc.NewAppConfig(context.Background(), tt.configPayload)
		require.Nil(t, err)

		payload := tt.configPayload
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, *payload.OrgInfo.OrgLogoURL, result.OrgLogoURL)
		assert.Equal(t, *payload.OrgInfo.OrgName, result.OrgName)
		assert.Equal(t, "https://acme.co:8080", result.KolideServerURL)
	}
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	var configInfo appConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&configInfo)
	require.Nil(t, err)
	require.NotNil(t, configInfo.SMTPSettings)
	config := configInfo.SMTPSettings
	assert.Equal(t, uint(465), *config.SMTPPort)
	require.NotNil(t, *configInfo.OrgInfo)
	assert.Equal(t, "Kolide", *configInfo.OrgInfo.OrgName)
	assert.Equal(t, "http://foo.bar/image.png", *configInfo.OrgInfo.OrgLogoURL)

}

func testModifyAppConfig(t *testing.T, r *testResource) {
	config := &kolide.AppConfig{
		KolideServerURL:        "https://foo.com",
		OrgName:                "Zip",
		OrgLogoURL:             "http://foo.bar/image.png",
		SMTPPort:               567,
		SMTPAuthenticationType: kolide.AuthTypeNone,
		SMTPServer:             "foo.com",
		SMTPEnableTLS:          true,
		SMTPVerifySSLCerts:     true,
		SMTPEnableStartTLS:     true,
		EnableSSO:              true,
		IDPName:                "idpname",
		Metadata:               "metadataxxxxxx",
		IssuerURI:              "http://issuer.idp.com",
		EntityID:               "kolide",
	}
	payload := appConfigPayloadFromAppConfig(config)
	payload.SMTPTest = new(bool)

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(payload)
	require.Nil(t, err)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/config", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var respBody appConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.Nil(t, err)
	require.NotNil(t, respBody.OrgInfo)
	assert.Equal(t, config.OrgName, *respBody.OrgInfo.OrgName)
	saved, err := r.ds.AppConfig()
	require.Nil(t, err)
	// verify email test succeeded
	assert.True(t, saved.SMTPConfigured)
	// verify that SSO stuff was saved
	assert.True(t, saved.EnableSSO)
	assert.Equal(t, "idpname", saved.IDPName)
	assert.Equal(t, "metadataxxxxxx", saved.Metadata)
	assert.Equal(t, "http://issuer.idp.com", saved.IssuerURI)
	assert.Equal(t, "kolide", saved.EntityID)

}

func testModifyAppConfigWithValidationFail(t *testing.T, r *testResource) {
	config := &kolide.AppConfig{
		SMTPEnableStartTLS: false,
	}
	payload := appConfigPayloadFromAppConfig(config)
	payload.SMTPTest = new(bool)

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(payload)
	require.Nil(t, err)
	req, err := http.NewRequest("PATCH", r.server.URL+"/api/v1/kolide/config", &buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var validationErrors mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&validationErrors)
	require.Nil(t, err)
	require.Equal(t, 0, len(validationErrors.Errors))
	existing, err := r.ds.AppConfig()
	assert.Nil(t, err)
	assert.Equal(t, config.SMTPEnableStartTLS, existing.SMTPEnableStartTLS)
}

func appConfigPayloadFromAppConfig(config *kolide.AppConfig) *kolide.AppConfigPayload {
	return &kolide.AppConfigPayload{
		OrgInfo: &kolide.OrgInfo{
			OrgLogoURL: &config.OrgLogoURL,
			OrgName:    &config.OrgName,
		},
		ServerSettings: &kolide.ServerSettings{
			KolideServerURL: &config.KolideServerURL,
		},
		SMTPSettings: smtpSettingsFromAppConfig(config),
		SSOSettings: &kolide.SSOSettingsPayload{
			EnableSSO:   &config.EnableSSO,
			IDPName:     &config.IDPName,
			Metadata:    &config.Metadata,
			MetadataURL: &config.MetadataURL,
			IssuerURI:   &config.IssuerURI,
			EntityID:    &config.EntityID,
		},
	}
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func TestSSONotPresent(t *testing.T) {
	invalid := &invalidArgumentError{}
	var p kolide.AppConfigPayload
	validateSSOSettings(p, &kolide.AppConfig{}, invalid)
	assert.False(t, invalid.HasErrors())

}

func TestNeedFieldsPresent(t *testing.T) {
	invalid := &invalidArgumentError{}
	config := kolide.AppConfig{
		EnableSSO:   true,
		EntityID:    "kolide",
		IssuerURI:   "http://issuer.idp.com",
		MetadataURL: "http://isser.metadata.com",
		IDPName:     "onelogin",
	}
	p := appConfigPayloadFromAppConfig(&config)
	validateSSOSettings(*p, &kolide.AppConfig{}, invalid)
	assert.False(t, invalid.HasErrors())
}

func TestMissingMetadata(t *testing.T) {
	invalid := invalidArgumentError{}
	config := kolide.AppConfig{
		EnableSSO: true,
		EntityID:  "kolide",
		IssuerURI: "http://issuer.idp.com",
		IDPName:   "onelogin",
	}
	p := appConfigPayloadFromAppConfig(&config)
	validateSSOSettings(*p, &kolide.AppConfig{}, &invalid)
	require.True(t, invalid.HasErrors())
	require.Len(t, invalid, 1)
	assert.Equal(t, "metadata", invalid[0].name)
	assert.Equal(t, "either metadata or metadata_url must be defined", invalid[0].reason)
}
