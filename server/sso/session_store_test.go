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
	store := NewSessionStore(p)
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

	err = store.CreateSession("ssohandle2", "http://bar.com", 1)
	require.Nil(t, err)
	sess, err := store.GetSession("ssohandle2")
	require.Nil(t, err)
	require.NotNil(t, sess)
	assert.Equal(t, "http://bar.com", sess.OriginalURL)
	sess, err = store.UpdateSession("ssohandle2", "user@xxx.com")
	require.Nil(t, err)
	assert.Equal(t, "http://bar.com", sess.OriginalURL)
	sess, err = store.GetSession("ssohandle2")
	require.Nil(t, err)
	assert.Equal(t, "user@xxx.com", sess.UserName)

	err = store.ExpireSession("ssohandle2")
	require.Nil(t, err)
	_, err = store.GetSession("ssohandle2")
	require.NotNil(t, err)
	assert.Equal(t, ErrSessionNotFound, err)

	err = store.CreateSession("ssohandle3", "", 5)
	require.Nil(t, err)
	err = store.CreateSession("ssohandle3", "", 5)
	require.NotNil(t, err)
	assert.Equal(t, "session key 'ssohandle3' is already in use", err.Error())

	err = store.CreateSession("short", "http://foo.com", 2)
	require.Equal(t, "session handle must be 8 or more characters in length", err.Error())
}

func TestEncrypt(t *testing.T) {
	if _, ok := os.LookupEnv("REDIS_TEST"); !ok {
		t.Skip("skipping sso session store tests")
	}
	pool := newPool(t)
	defer pool.Close()
	st := &store{pool: pool}
	key := CreateAESKey()
	ssoHandle := "thistest"
	relayState, err := st.EncryptSSOHandle(ssoHandle, key, 2)
	require.Nil(t, err)
	assert.NotEqual(t, ssoHandle, relayState)

	ssoHandle2, err := st.DescryptSSOHandle(relayState, key)
	require.Nil(t, err)
	assert.Equal(t, ssoHandle, ssoHandle2)
}
