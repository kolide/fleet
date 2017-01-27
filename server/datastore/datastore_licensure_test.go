package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLicensure(t *testing.T, ds kolide.Datastore) {
	if ds.Name() == "inmem" {
		t.Skip("inmem is being deprecated")
	}
	err := ds.MigrateData()
	require.Nil(t, err)
	license, err := ds.License()
	require.Nil(t, err)
	assert.Nil(t, license.TokenString)

	err = ds.SaveLicense("fake license")
	require.Nil(t, err)
	license, err = ds.License()
	require.Nil(t, err)
	require.NotNil(t, license.TokenString)
	assert.Equal(t, "fake license", *license.TokenString)

}
