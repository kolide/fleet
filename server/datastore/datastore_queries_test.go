package datastore

import (
	"fmt"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/patrickmn/sortutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDeleteQuery(t *testing.T, ds kolide.Datastore) {
	query := &kolide.Query{
		Name:     "foo",
		Query:    "bar",
		Interval: 123,
	}
	query, err := ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotEqual(t, query.ID, 0)

	err = ds.DeleteQuery(query)
	assert.Nil(t, err)

	assert.NotEqual(t, query.ID, 0)
	_, err = ds.Query(query.ID)
	assert.NotNil(t, err)
}

func testSaveQuery(t *testing.T, ds kolide.Datastore) {
	query := &kolide.Query{
		Name:  "foo",
		Query: "bar",
	}
	query, err := ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, query.ID)

	query.Query = "baz"
	err = ds.SaveQuery(query)

	assert.Nil(t, err)

	queryVerify, err := ds.Query(query.ID)
	assert.Nil(t, err)
	assert.Equal(t, "baz", queryVerify.Query)
}

func testListQuery(t *testing.T, ds kolide.Datastore) {
	for i := 0; i < 10; i++ {
		_, err := ds.NewQuery(&kolide.Query{
			Name:  fmt.Sprintf("name%02d", i),
			Query: fmt.Sprintf("query%02d", i),
		})
		assert.Nil(t, err)
	}

	opts := kolide.ListOptions{}
	results, err := ds.ListQueries(opts)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(results))
}

func testDistributedQueryCampaign(t *testing.T, ds kolide.Datastore) {
	mockClock := clock.NewMockClock()

	query, err := ds.NewQuery(&kolide.Query{
		Name:  "test",
		Query: "select * from time",
	})
	require.Nil(t, err)

	campaign, err := ds.NewDistributedQueryCampaign(&kolide.DistributedQueryCampaign{
		QueryID: query.ID,
		Status:  kolide.QueryRunning,
	})
	require.Nil(t, err)

	{
		retrieved, err := ds.DistributedQueryCampaign(campaign.ID)
		require.Nil(t, err)
		assert.Equal(t, campaign.QueryID, retrieved.QueryID)
		assert.Equal(t, campaign.Status, retrieved.Status)
	}

	h1, err := ds.NewHost(&kolide.Host{
		HostName:         "foo.local",
		PrimaryIP:        "192.168.1.10",
		NodeKey:          "1",
		UUID:             "1",
		DetailUpdateTime: mockClock.Now(),
	})
	require.Nil(t, err)
	require.Nil(t, ds.MarkHostSeen(h1, mockClock.Now()))

	h2, err := ds.NewHost(&kolide.Host{
		HostName:         "bar.local",
		PrimaryIP:        "192.168.1.11",
		NodeKey:          "2",
		UUID:             "2",
		DetailUpdateTime: mockClock.Now().Add(-1 * time.Hour),
	})
	require.Nil(t, err)
	// make this host "offline"
	require.Nil(t, ds.MarkHostSeen(h2, mockClock.Now().Add(-1*time.Hour)))

	h3, err := ds.NewHost(&kolide.Host{
		HostName:         "baz.local",
		PrimaryIP:        "192.168.1.12",
		NodeKey:          "3",
		UUID:             "3",
		DetailUpdateTime: mockClock.Now().Add(-13 * time.Minute),
	})
	require.Nil(t, err)
	require.Nil(t, ds.MarkHostSeen(h3, mockClock.Now().Add(-5*time.Minute)))

	l1, err := ds.NewLabel(&kolide.Label{
		Name:  "label foo",
		Query: "query foo",
	})
	require.Nil(t, err)
	require.NotZero(t, l1.ID)

	l2, err := ds.NewLabel(&kolide.Label{
		Name:  "label bar",
		Query: "query foo",
	})
	require.Nil(t, err)
	require.NotZero(t, l2.ID)

	addHost := func(h *kolide.Host) {
		_, err := ds.NewDistributedQueryCampaignTarget(
			&kolide.DistributedQueryCampaignTarget{
				Type:                       kolide.TargetHost,
				TargetID:                   h.ID,
				DistributedQueryCampaignID: campaign.ID,
			})
		require.Nil(t, err)

	}

	addLabel := func(l *kolide.Label) {
		_, err := ds.NewDistributedQueryCampaignTarget(
			&kolide.DistributedQueryCampaignTarget{
				Type:                       kolide.TargetLabel,
				TargetID:                   l.ID,
				DistributedQueryCampaignID: campaign.ID,
			})
		require.Nil(t, err)
	}

	checkTargets := func(expectedHostIDs []uint, expectedLabelIDs []uint) {
		hostIDs, labelIDs, err := ds.DistributedQueryCampaignTargetIDs(campaign.ID)
		require.Nil(t, err)

		sortutil.Asc(expectedHostIDs)
		sortutil.Asc(hostIDs)
		assert.Equal(t, expectedHostIDs, hostIDs)

		sortutil.Asc(expectedLabelIDs)
		sortutil.Asc(labelIDs)
		assert.Equal(t, expectedLabelIDs, labelIDs)
	}

	checkTargets([]uint{}, []uint{})

	addHost(h1)
	checkTargets([]uint{h1.ID}, []uint{})

	addLabel(l1)
	checkTargets([]uint{h1.ID}, []uint{l1.ID})

	addLabel(l2)
	checkTargets([]uint{h1.ID}, []uint{l1.ID, l2.ID})

	addHost(h2)
	addHost(h3)
	checkTargets([]uint{h1.ID, h2.ID, h3.ID}, []uint{l1.ID, l2.ID})

}
