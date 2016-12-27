package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/datastore/internal/appstate"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testOptions(t *testing.T, ds kolide.Datastore) {
	// were options pre-loaded?
	opts, err := ds.Options()
	require.Nil(t, err)
	assert.Len(t, opts, len(appstate.Options))

	opt, err := ds.OptionByName("aws_access_key_id")
	require.Nil(t, err)
	require.NotNil(t, opt)
	opt2, err := ds.Option(opt.ID)
	require.Nil(t, err)
	require.NotNil(t, opt2)
	assert.Equal(t, opt.Name, opt2.Name)
	assert.Equal(t, opt.Value, opt2.Value)

	opt.Value = new(string)
	*opt.Value = "true"
	err = ds.SaveOptions([]kolide.Option{*opt})
	require.Nil(t, err)

	// can't change a read only option
	opt, err = ds.OptionByName("disable_distributed")
	require.Nil(t, err)
	opt.Value = new(string)
	*opt.Value = "true"
	err = ds.SaveOptions([]kolide.Option{*opt})
	require.NotNil(t, err)

	opt, _ = ds.OptionByName("aws_profile_name")
	assert.Nil(t, opt.Value)
	opt.Value = new(string)
	*opt.Value = "zip"
	opt2, _ = ds.OptionByName("disable_distributed")
	assert.Equal(t, "false", *opt2.Value)
	*opt2.Value = "true"
	modList := []kolide.Option{*opt, *opt2}
	// The aws access key option can be saved but because the disable_events can't
	// be we want to verify that the whole transaction is rolled back
	err = ds.SaveOptions(modList)
	assert.NotNil(t, err)
	opt, _ = ds.OptionByName("aws_profile_name")
	assert.Nil(t, opt.Value)
	opt2, err = ds.OptionByName("disable_distributed")
	assert.Equal(t, "false", *opt2.Value)

}
