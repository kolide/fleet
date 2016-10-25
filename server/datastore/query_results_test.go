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
	// Test handling results for two campaigns in parallel

	campaign1 := kolide.DistributedQueryCampaign{ID: 1}

	channel1, err := store.ReadChannel(campaign1)
	assert.Nil(t, err)

	results1 := []kolide.DistributedQueryResult{}

	expected1 := []kolide.DistributedQueryResult{
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

	campaign2 := kolide.DistributedQueryCampaign{ID: 2}

	channel2, err := store.ReadChannel(campaign2)
	assert.Nil(t, err)

	results2 := []kolide.DistributedQueryResult{}

	expected2 := []kolide.DistributedQueryResult{
		kolide.DistributedQueryResult{
			DistributedQueryCampaignID: 2,
			ResultJSON:                 json.RawMessage(`{"tim":"tom"}`),
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
			DistributedQueryCampaignID: 2,
			ResultJSON:                 json.RawMessage(`{"slim":"slam"}`),
			Host: kolide.Host{
				ID:               3,
				UpdatedAt:        time.Now(),
				DetailUpdateTime: time.Now(),
			},
		},
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for res := range channel1 {
			results1 = append(results1, res)
		}

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for res := range channel2 {
			results2 = append(results2, res)
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, res := range expected1 {
			assert.Nil(t, store.WriteResult(res))
		}
		store.CloseQuery(campaign1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, res := range expected2 {
			assert.Nil(t, store.WriteResult(res))
		}
		store.CloseQuery(campaign2)
	}()

	// wait with a timeout to ensure that the test can't hang
	if waitTimeout(&wg, 2*time.Second) {
		t.Error("Timed out waiting for goroutines to join")
	}

	assert.EqualValues(t, expected1, results1)
	assert.EqualValues(t, expected2, results2)

}
