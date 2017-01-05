package service

import (
	"testing"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/datastore/inmem"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createServiceMockForImport(t *testing.T) *service {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)
	return &service{
		ds: ds,
	}
}

func TestOptionsImportConfig(t *testing.T) {
	opts := kolide.OptionNameToValueMap{
		"aws_access_key_id": "foo",
	}
	resp := kolide.NewImportConfigResponse()
	svc := createServiceMockForImport(t)
	err := svc.importOptions(opts, resp)
	require.Nil(t, err)
	status := resp.Status(kolide.OptionsSection)
	require.NotNil(t, status)
	assert.Equal(t, 1, status.ImportCount)
	opt, err := svc.ds.OptionByName("aws_access_key_id")
	require.Nil(t, err)
	assert.Equal(t, "foo", opt.GetValue())
	require.Len(t, status.Messages, 1)
	assert.Equal(t, "set aws_access_key_id value to foo", status.Messages[0])
}

func TestOptionsImportConfigWithSkips(t *testing.T) {
	opts := kolide.OptionNameToValueMap{
		"aws_access_key_id":     "foo",
		"aws_secret_access_key": "secret",
		// this should be skipped because it's already set
		"aws_firehose_period": 100,
		// these should be skipped because it's read only
		"disable_distributed": false,
		"pack_delimiter":      "x",
		// this should be skipped because it's not an option we know about
		"wombat": "not venomous",
	}
	resp := kolide.NewImportConfigResponse()
	svc := createServiceMockForImport(t)
	// set option val, it should be skipped
	opt, err := svc.ds.OptionByName("aws_firehose_period")
	require.Nil(t, err)
	opt.SetValue(23)
	err = svc.ds.SaveOptions([]kolide.Option{*opt})
	require.Nil(t, err)
	err = svc.importOptions(opts, resp)
	require.Nil(t, err)
	status := resp.Status(kolide.OptionsSection)
	require.NotNil(t, status)
	assert.Equal(t, 2, status.ImportCount)
	assert.Equal(t, 4, status.SkipCount)
	assert.Len(t, status.Warnings[kolide.OptionAlreadySet], 1)
	assert.Len(t, status.Warnings[kolide.OptionReadonly], 2)
	assert.Len(t, status.Warnings[kolide.OptionUnknown], 1)
	assert.Len(t, status.Messages, 2)
}
