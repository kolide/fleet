package sso

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kolide/kolide/server/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPool(t *testing.T) *redis.Pool {
	if _, ok := os.LookupEnv("REDIS_TEST"); ok {
		var (
			addr     = "127.0.0.1:6379"
			password = ""
		)
		if a, ok := os.LookupEnv("REDIS_PORT_6379_TCP_ADDR"); ok {
			addr = fmt.Sprintf("%s:6379", a)
		}

		p := pubsub.NewRedisPool(addr, password)
		_, err := p.Get().Do("PING")
		require.Nil(t, err)
		return p
	}
	return nil
}

func TestSessionStore(t *testing.T) {
	if _, ok := os.LookupEnv("REDIS_TEST"); !ok {
		t.Skip("skipping sso session store tests")
	}
	p := newPool(t)
	require.NotNil(t, p)
	defer p.Close()
	store := NewSessionStorage(p)
	require.NotNil(t, store)
	// make sure session expires, create session that 'lives'
	// for a second
	err := store.CreateSession("ssohandle", "http://foo.com", 1)
	require.Nil(t, err)
	// wait 2 seconds
	time.Sleep(2 * time.Second)
	_, err = store.GetSession("ssohandle")
	require.NotNil(t, err)
	assert.Equal(t, ErrSessionNotFound, err)

	err = store.CreateSession("handle2", "http://bar.com", 30)
	require.Nil(t, err)
	sess, err := store.GetSession("handle2")
	require.Nil(t, err)
	require.NotNil(t, sess)
	assert.Equal(t, "http://bar.com", sess.OriginalURL)
	sess.UserName = "user@xxx.com"
	sess.OriginalURL = "http://new.com"
	err = store.UpdateSession("handle2", sess, 10)
	require.Nil(t, err)
	getSess, err := store.GetSession("handle2")
	require.Nil(t, err)
	require.NotNil(t, getSess)

	assert.Equal(t, "user@xxx.com", getSess.UserName)
	assert.Equal(t, "http://new.com", getSess.OriginalURL)

	err = store.ExpireSession("handle2")
	require.Nil(t, err)
	_, err = store.GetSession("handle2")
	require.NotNil(t, err)
	assert.Equal(t, ErrSessionNotFound, err)

	err = store.CreateSession("dup", "", 5)
	require.Nil(t, err)
	err = store.CreateSession("dup", "", 5)
	require.NotNil(t, err)
	assert.Equal(t, "session key 'dup' is already in use", err.Error())
}
