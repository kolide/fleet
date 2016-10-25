package datastore

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out. http://stackoverflow.com/a/32843750/491710
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func TestQueryResultsStore(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPass := os.Getenv("REDIS_PASSWORD")
	if redisAddr != "" {
		store := newRedisQueryResults(redisAddr, redisPass)
		t.Run("redis", func(t *testing.T) {
			_, err := store.pool.Get().Do("PING")
			require.Nil(t, err)
			testQueryResultsStore(t, &store)
		})
	} else {
		t.Log("Skipping redis")
	}

	store := newInmemQueryResults()
	t.Run("inmem", func(t *testing.T) {
		testQueryResultsStore(t, &store)
	})
}

func testQueryResultsStore(t *testing.T, store kolide.QueryResultStore) {

	campaign := kolide.DistributedQueryCampaign{ID: 1}

	channel, err := store.ReadChannel(campaign)
	assert.Nil(t, err)

	results := []kolide.DistributedQueryResult{}

	expected := []kolide.DistributedQueryResult{
		kolide.DistributedQueryResult{
			DistributedQueryCampaignID: 1,
			ResultJSON:                 json.RawMessage(`{"foo":"bar"}`),
			Host: kolide.Host{
				ID: 1,
				// Note these times need to be set to avoid
				// issues with roundtrip serializing the zero
				// time value. See https://goo.gl/CCEs8x
				UpdatedAt:        time.Now(),
				DetailUpdateTime: time.Now(),
			},
		},
		kolide.DistributedQueryResult{
			DistributedQueryCampaignID: 1,
			ResultJSON:                 json.RawMessage(`{"whoo":"wahh"}`),
			Host: kolide.Host{
				ID:               3,
				UpdatedAt:        time.Now(),
				DetailUpdateTime: time.Now(),
			},
		},
		kolide.DistributedQueryResult{
			DistributedQueryCampaignID: 1,
			ResultJSON:                 json.RawMessage(`{"bing":"fds"}`),
			Host: kolide.Host{
				ID:               4,
				UpdatedAt:        time.Now(),
				DetailUpdateTime: time.Now(),
			},
		},
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for res := range channel {
			results = append(results, res)
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, res := range expected {
			assert.Nil(t, store.WriteResult(res))
		}
		store.CloseQuery(campaign)
	}()

	// wait with a timeout to ensure that the test can't hang
	if waitTimeout(&wg, 2*time.Second) {
		t.Error("Timed out waiting for goroutines to join")
	}

	assert.EqualValues(t, expected, results)

}
