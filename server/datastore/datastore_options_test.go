package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testOptions(t *testing.T, ds kolide.Datastore) {
	option, err := ds.NewOption("optString", kolide.OptionTypeString, false)
	require.Nil(t, err)
	assert.NotNil(t, option)
	opts, err := ds.Options()
	require.Nil(t, err)
	assert.Equal(t, 1, len(opts))

	_, err = ds.NewOption("optInt", kolide.OptionTypeInt, false)
	require.Nil(t, err)
	_, err = ds.NewOption("optFlag", kolide.OptionTypeFlag, false)
	require.Nil(t, err)
	// can't create dup option
	_, err = ds.NewOption("optFlag", kolide.OptionTypeFlag, false)
	assert.NotNil(t, err)
	opts, err = ds.Options()
	require.Nil(t, err)
	assert.Equal(t, 3, len(opts))
	optVals := []kolide.OptionValue{
		kolide.OptionValue{
			OptionID: opts[0].ID,
			Value:    "text",
		},
		kolide.OptionValue{
			OptionID: opts[1].ID,
			Value:    "23",
		},
	}
	result, err := ds.SetOptionValues(optVals)
	require.Nil(t, err)
	assert.Equal(t, 2, len(result))

	optVals = []kolide.OptionValue{
		kolide.OptionValue{
			OptionID: opts[0].ID,
			Value:    "text",
		},
	}
	_, err = ds.SetOptionValues(optVals)
	require.Nil(t, err)

	ovs, err := ds.OptionValues()
	require.Nil(t, err)
	assert.Equal(t, 1, len(ovs))

}
