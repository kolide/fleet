package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjZjOjQyOjNkOjc3OjhkOmE" +
	"4OjAyOjdjOmFjOjE2OmUxOmIwOjYzOjZjOmM1OjdiIn0.eyJsaWNlbnNlX3V1aWQi" +
	"OiIzNTA3YzgxMy0yNDRmLTRjOTEtYjJiOC0xZGRjNDU1ZjFlNzQiLCJvcmdhbml" +
	"6YXRpb25fbmFtZSI6IlBoYW50YXNtLCBJbmMuIiwib3JnYW5pemF0aW9uX3V1aWQ" +
	"iOiI4NmE5MTFhZS05MmEzLTRkOWUtYjEzZS0yYTZmZDc5N2ZlMTkiLCJob3N0X2xp" +
	"bWl0IjowLCJldmFsdWF0aW9uIjp0cnVlLCJleHBpcmVzX2F0IjoiMjAxNy0wMy0wMiA" +
	"yMzoxODoyMSBVVEMifQ.LOZ35ZwGqN_AJ_A0TDNKHfrHy4OCP9SM0NRzunUvU-qhPo_Td" +
	"f_omU0KLmN0aeWOlgBFVD-5kE3xLSqO2W0n1bU2ktJFZDpCe_yrKZK3e4noV4QZtHsjVKWObTu2" +
	"s8EfJNw8qMYvatv5AJ77Qbnnbf_Ic5eoP0_mLm-wHkXukl3PutR82dmIbwcR8gjf-ZZBYr21Q" +
	"QoFjIUGWqH1ttN__sRoTrMFI1XKSfH-5CvuIi69iJYiymZwCUGbDOdyGRVOIij-9xfZP15s" +
	"d8jseT_yRjpa4FF55UmGozkAAkTUnGFbC-RJWhiWNFBZbM26xlKWRSsAR6h1pQ9Qmdiihtnw4w"

func testLicensure(t *testing.T, ds kolide.Datastore) {
	if ds.Name() == "inmem" {
		t.Skip("inmem is being deprecated")
	}
	err := ds.MigrateData()
	require.Nil(t, err)
	license, err := ds.License()
	require.Nil(t, err)
	assert.Nil(t, license.Token)

	publicKey, err := ds.PublicKey(token)
	require.Nil(t, err)

	_, err = ds.SaveLicense(token, publicKey)
	require.Nil(t, err)
	license, err = ds.License()
	require.Nil(t, err)
	require.NotNil(t, license.Token)
	assert.Equal(t, token, *license.Token)

}
