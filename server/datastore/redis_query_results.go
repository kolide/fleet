package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kolide/kolide-ose/server/kolide"
)

type redisQueryResults struct {
	// mapping of query ID to corresponding connection
	connections map[uint]redis.PubSubConn
	// connection pool
	pool redis.Pool
	// mutex used to protect the redis connection struct from race
	// conditions during unsubscribe
	mutex sync.Mutex
}

var _ kolide.QueryResultStore = &redisQueryResults{}

func newRedisQueryResults(server, password string) redisQueryResults {
	return redisQueryResults{
		connections: map[uint]redis.PubSubConn{},
		pool: redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func channelForID(id uint) string {
	return fmt.Sprintf("results_%d", id)
}

func (im *redisQueryResults) WriteResult(result kolide.DistributedQueryResult) error {
	conn := im.pool.Get()
	defer conn.Close()

	channelName := channelForID(result.DistributedQueryCampaignID)

	jsonVal, err := json.Marshal(&result)
	if err != nil {
		return errors.New("error marshalling JSON for writing result: " + err.Error())
	}

	n, err := redis.Int(conn.Do("PUBLISH", channelName, string(jsonVal)))
	if err != nil {
		return fmt.Errorf("PUBLISH failed to channel %s: %s", channelName, err.Error())
	}
	if n == 0 {
		return fmt.Errorf("no subscribers for channel %s", channelName)
	}

	return nil
}

func (im *redisQueryResults) ReadChannel(query kolide.DistributedQueryCampaign) (<-chan kolide.DistributedQueryResult, error) {
	if _, exists := im.connections[query.ID]; exists {
		return nil, fmt.Errorf("channel already open for ID %d", query.ID)
	}

	outChannel := make(chan kolide.DistributedQueryResult)

	conn := redis.PubSubConn{Conn: im.pool.Get()}

	channelName := channelForID(query.ID)
	conn.Subscribe(channelName)

	// Save connection so it can be unsubscribed from CloseQuery
	// and therefore alert us through the Receive() method below
	im.connections[query.ID] = conn

	go func() {
		defer func() {
			im.mutex.Lock()
			close(outChannel)
			conn.Unsubscribe()
			conn.Close()
			im.mutex.Unlock()
		}()

		for {
			switch v := conn.Receive().(type) {

			case redis.Message:
				var res kolide.DistributedQueryResult
				err := json.Unmarshal(v.Data, &res)
				if err != nil {
					fmt.Println(err)
				}
				outChannel <- res

			case redis.Subscription:
				if v.Channel == channelName && v.Count == 0 {
					return
				}

			case error:
				return
			}
		}

	}()
	return outChannel, nil
}

func (im *redisQueryResults) CloseQuery(query kolide.DistributedQueryCampaign) {
	conn, ok := im.connections[query.ID]
	if !ok {
		return
	}
	im.mutex.Lock()
	conn.Unsubscribe()
	im.mutex.Unlock()
}
