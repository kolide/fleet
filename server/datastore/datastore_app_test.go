package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testOrgInfo(t *testing.T, ds kolide.Datastore) {
	info := &kolide.AppConfig{
		OrgName:    "Kolide",
		OrgLogoURL: "localhost:8080/logo.png",
		SMTPConfig: &kolide.SMTPConfig{},
	}

	info, err := ds.NewAppConfig(info)
	assert.Nil(t, err)

	info2, err := ds.AppConfig()
	require.Nil(t, err)
	assert.Equal(t, info2.OrgName, info.OrgName)
	assert.False(t, info2.Configured)

	info2.OrgName = "koolide"
	info2.Domain = "foo"
	info2.Configured = true
	info2.SenderAddress = "123"
	info2.Server = "server"
	info2.Port = 100
	info2.AuthenticationType = kolide.AuthTypeUserNamePassword
	info2.UserName = "username"
	info2.Password = "password"
	info2.EnableTLS = false
	info2.AuthenticationMethod = kolide.AuthMethodCramMD5
	info2.VerifySSLCerts = true
	info2.EnableStartTLS = true
	err = ds.SaveAppConfig(info2)
	require.Nil(t, err)

	info3, err := ds.AppConfig()
	require.Nil(t, err)
	assert.Equal(t, info2, info3)

	info4, err := ds.NewAppConfig(info3)
	assert.Nil(t, err)
	assert.Equal(t, info3, info4)
}
