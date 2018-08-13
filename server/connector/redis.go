package connector

import (
	"github.com/gomodule/redigo/redis"
	//"github.com/kolide/fleet/server/health"
	"github.com/kolide/fleet/server/config"
	"time"
)

// NewRedisPool creates a Redis connection pool using the provided server
// address and password.
func NewRedisConn(conf config.RedisConfig) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Address)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
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
	}, nil 
}

type redisHealthChecker struct {
	// connection pool
	conn *redis.Pool
}

//var _ health.Checker = &redisHealthChecker{}

func NewRedisHealthChecker(conn *redis.Pool) (*redisHealthChecker, error) {
	return &redisHealthChecker{conn: conn}, nil 
}

// HealthCheck verifies that the redis backend can be pinged, returning an error
// otherwise.
func (r *redisHealthChecker) HealthCheck() error {
	conn := r.conn.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err
}
