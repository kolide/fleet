package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testOptions(t *testing.T, ds kolide.Datastore) {
	opt, err := kolide.NewOption("name", "option", kolide.OptionTypeString, true)
	require.Nil(t, err)
	require.NotNil(t, opt)
	opt2, err := ds.SaveOption(*opt)
	require.Nil(t, err)
	assert.NotEqual(t, 0, opt2.ID)
	assert.Equal(t, opt.Name, opt2.Name)
	assert.Equal(t, opt.RawValue, opt2.RawValue)

	_, err = kolide.NewOption("name2", 3, kolide.OptionTypeString, true)
	require.NotNil(t, err)
	assert.Equal(t, "type mismatch", err.Error())

	opt3, err := kolide.NewOption("name3", nil, kolide.OptionTypeFlag, false)
	opt4, err := ds.SaveOption(*opt3)
	require.Nil(t, err)
	assert.Nil(t, opt4.RawValue)

}

func testOptionsSetup(t *testing.T, ds kolide.Datastore) {
	err := kolide.CreateOptions(ds)
	require.Nil(t, err)
}
