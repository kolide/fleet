package service

import (
	"testing"

	"github.com/WatchBeam/clock"
	"github.com/kolide/kolide/server/config"
	"github.com/kolide/kolide/server/datastore/mysql"
	"github.com/kolide/kolide/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestCreateAppConfig(t *testing.T) {
	ds, err := mysql.New(config.TestConfig().Mysql, clock.NewMockClock())
	require.Nil(t, err)
	require.Nil(t, ds.MigrateData())

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)
	var appConfigTests = []struct {
		configPayload kolide.AppConfigPayload
	}{
		{
			configPayload: kolide.AppConfigPayload{
				OrgInfo: &kolide.OrgInfo{
					OrgLogoURL: stringPtr("acme.co/images/logo.png"),
					OrgName:    stringPtr("Acme"),
				},
				ServerSettings: &kolide.ServerSettings{
					KolideServerURL: stringPtr("https://acme.co:8080/"),
				},
			},
		},
	}

	for _, tt := range appConfigTests {
		result, err := svc.NewAppConfig(context.Background(), tt.configPayload)
		require.Nil(t, err)

		payload := tt.configPayload
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, *payload.OrgInfo.OrgLogoURL, result.OrgLogoURL)
		assert.Equal(t, *payload.OrgInfo.OrgName, result.OrgName)
		assert.Equal(t, *payload.ServerSettings.KolideServerURL, result.KolideServerURL)
	}
}
