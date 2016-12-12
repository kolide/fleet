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
	info2.EnableSSLTLS = false
	info2.AuthenticationMethod = kolide.AuthMethodCramMD5
	info2.VerifySSLCerts = true
	info2.EnableStartTLS = true
	err = ds.SaveAppConfig(info2)
	require.Nil(t, err)

	info3, err := ds.AppConfig()
	require.Nil(t, err)
	assert.Equal(t, info3.OrgName, info2.OrgName)
	assert.Equal(t, info3.Domain, info2.Domain)
	assert.Equal(t, info3.Configured, info2.Configured)
	assert.Equal(t, info3.SenderAddress, info2.SenderAddress)
	assert.Equal(t, info3.Server, info2.Server)
	assert.Equal(t, info3.Port, info2.Port)
	assert.Equal(t, info3.AuthenticationType, info2.AuthenticationType)
	assert.Equal(t, info3.UserName, info2.UserName)
	assert.Equal(t, info3.Password, info2.Password)
	assert.Equal(t, info3.EnableSSLTLS, info2.EnableSSLTLS)
	assert.Equal(t, info3.AuthenticationMethod, info2.AuthenticationMethod)
	assert.Equal(t, info3.VerifySSLCerts, info2.VerifySSLCerts)
	assert.Equal(t, info3.EnableStartTLS, info2.EnableStartTLS)

	info4, err := ds.NewAppConfig(info3)
	assert.Nil(t, err)
	assert.Equal(t, info4.OrgName, info3.OrgName)
}
