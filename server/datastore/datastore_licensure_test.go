package datastore

import (
	"testing"

	"github.com/kolide/kolide/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjRkOmM1OmRlOmE1OjczOmUxOmE4OjI4OmU2OmEyOjMwOmI4OmI1OjBmOjg4OjQ0In0.eyJsaWNlbnNlX3V1aWQiOiIyYWYyZDlhMC1iOWE1LTQ0ZTItODU1NC04Mjc2MGI4ODQwZDYiLCJvcmdhbml6YXRpb25fbmFtZSI6IlBoYW50YXNtLCBJbmMuIiwib3JnYW5pemF0aW9uX3V1aWQiOiJkZmJkNWIwMy0xMDg0LTQ2YWUtYjM4MS1lMTI5YWM2NmU4ZDgiLCJob3N0X2xpbWl0IjowLCJldmFsdWF0aW9uIjp0cnVlLCJleHBpcmVzX2F0IjoiMjAxNy0wMy0wNFQxNTowMTo0OSswMDowMCJ9.Ny4Fxqlq_4U647gmIouFPZQH4YG8R_AHOlDTBObWOUfhcKiz44vRkCqr_Jqprb0zVtSVy1bMojLLmQhKjSxQZuiqvQBfou9Osfd5D3i-TXEb5JpoCgFem-1t5jvOT7T9H4HJpuKE40cnOl3Zu2OzjjdxMMZbj_i2iwZytW1b7SrGNAwJVXXwJs2a95bGbMuZWyV-YpuHaWlx-VpTv4c2vQo2eQWTpTH7YdcQ7Mo_5QdN7247qKo_ORTtqLLTjg7BoxB__ydWMhxOQuRJGQAMc0OsZ72uLd7JKzvWpSLFk7mdVk718mweq6X2R0BPKtTc6lYjbPScoTysM2Owe5Hi7A"

func testLicensure(t *testing.T, ds kolide.Datastore) {
	if ds.Name() == "inmem" {
		t.Skip("inmem is being deprecated")
	}
	err := ds.MigrateData()
	require.Nil(t, err)
	license, err := ds.License()
	require.Nil(t, err)
	assert.Nil(t, license.Token)

	publicKey, err := ds.LicensePublicKey(token)
	require.Nil(t, err)

	_, err = ds.SaveLicense(token, publicKey)
	require.Nil(t, err)
	license, err = ds.License()
	require.Nil(t, err)
	require.NotNil(t, license.Token)
	assert.Equal(t, token, *license.Token)

}
