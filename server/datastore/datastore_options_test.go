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
	assert.Equal(t, opt.RawValue, opt2.RawValue)

	opt.RawValue = new(string)
	*opt.RawValue = "true"
	err = ds.SaveOption(*opt)
	require.Nil(t, err)

	opt, err = ds.OptionByName("disable_events")
	require.Nil(t, err)
	opt.RawValue = new(string)
	*opt.RawValue = "true"
	err = ds.SaveOption(*opt)
	require.NotNil(t, err)
	assert.Equal(t, "readonly option can't be changed", err.Error())
}
