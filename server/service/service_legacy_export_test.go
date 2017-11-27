package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportConfig(t *testing.T) {
	ds := &mock.Store{
		AppConfigStore: mock.AppConfigStore{
			AppConfigFunc: func() (*kolide.AppConfig, error) {
				return &kolide.AppConfig{
					FIMInterval: 300,
				}, nil
			},
		},
		OptionStore: mock.OptionStore{
			GetOsqueryConfigOptionsFunc: func() (map[string]interface{}, error) {
				return map[string]interface{}{
					"disable_distributed":           false,
					"distributed_interval":          10,
					"distributed_tls_read_endpoint": "/api/v1/osquery/distributed/read",
				}, nil
			},
		},
		FileIntegrityMonitoringStore: mock.FileIntegrityMonitoringStore{
			FIMSectionsFunc: func() (kolide.FIMSections, error) {
				return kolide.FIMSections{
					"etc": []string{
						"/etc/config/%%",
						"/etc/zipp",
					},
				}, nil
			},
		},
	}
	svc := service{
		ds: ds,
	}
	resp, err := svc.ExportConfig(context.Background())
	require.Nil(t, err)

	var o kolide.OptionsYaml
	err = yaml.Unmarshal([]byte(resp), &o)
	assert.Nil(t, err)
	assert.Equal(t, "TODO", o.ApiVersion)
	assert.Equal(t, "TODO", o.Kind)
	assert.JSONEq(t, `
		{
		  "file_paths":{
		    "etc":[
		      "/etc/config/%%",
		      "/etc/zipp"
		    ]
		  },
		  "options":{
		    "disable_distributed":false,
		    "distributed_interval":10,
		    "distributed_tls_read_endpoint":"/api/v1/osquery/distributed/read"
		  }
		}
		`,
		string(o.Spec.Config),
	)
}

func TestExportConfigErrors(t *testing.T) {
	var appConfigErr, optionErr, FIMErr error
	ds := &mock.Store{
		AppConfigStore: mock.AppConfigStore{
			AppConfigFunc: func() (*kolide.AppConfig, error) {
				if appConfigErr != nil {
					return nil, appConfigErr
				}
				return &kolide.AppConfig{}, nil
			},
		},
		OptionStore: mock.OptionStore{
			GetOsqueryConfigOptionsFunc: func() (map[string]interface{}, error) {
				if optionErr != nil {
					return nil, optionErr
				}
				return map[string]interface{}{}, nil
			},
		},
		FileIntegrityMonitoringStore: mock.FileIntegrityMonitoringStore{
			FIMSectionsFunc: func() (kolide.FIMSections, error) {
				if FIMErr != nil {
					return nil, FIMErr
				}
				return kolide.FIMSections{}, nil
			},
		},
	}
	svc := service{
		ds: ds,
	}
	_, err := svc.ExportConfig(context.Background())
	require.Nil(t, err)

	appConfigErr, optionErr, FIMErr = errors.New("foobar"), nil, nil
	_, err = svc.ExportConfig(context.Background())
	require.NotNil(t, err)

	appConfigErr, optionErr, FIMErr = nil, errors.New("foobar"), nil
	_, err = svc.ExportConfig(context.Background())
	require.NotNil(t, err)

	appConfigErr, optionErr, FIMErr = nil, nil, errors.New("foobar")
	_, err = svc.ExportConfig(context.Background())
	require.NotNil(t, err)
}
