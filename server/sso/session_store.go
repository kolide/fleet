package sso

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// Session stores state for the lifetime of a single sign on session
type Session struct {
	// OriginalURL is the resource being accessed when login request was triggered
	OriginalURL string `json:"original_url"`
	// UserName is only assigned from the IDP auth response, if present it
	// indicates the that user has authenticated against the IDP.
	UserName string `json:"user_name"`
	// ExpiresAt session will be removed after this time.
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionStore persists state of a sso session across process boundries and
// method calls by associating the state of the sign on session with a unique
// token created by the user agent (browser SPA).  The lifetime of the state object
// is constrained in the backing store (Redis) so if the sso process is not completed in
// a reasonable amount of time, it automatically expires and is removed.
type SessionStore interface {
	// CreateSession creates a sso session.
	// The ssoHandle is a unique token created by the
	// user agent. The originalURL maps to the resource that the user agent was
	// attempting to access when the auth challange was issued. lifetime is the amount of
	// time in seconds that the session will live.
	CreateSession(ssoHandle, originalURL string, lifetime uint) error
	// GetSession returns a session identfied by ssoHandle.
	GetSession(ssoHandle string) (*Session, error)
	// UpdateSession adds a user name to an existing sso session.
	UpdateSession(ssoHandle, userName string) (*Session, error)
	// ExpireSession removes session information from Redis
	ExpireSession(ssoHandle string) error
	// EncryptSSOHandle obfuscates the ssoHandle so it can't be used by external agents to
	// log in. We store the initialization vector we user to encrypt in redis so the
	// lifetime parameter makes sure that the record gets cleaned up. The key is the
	// AES encryption key.
	EncryptSSOHandle(ssoHandle string, key []byte, lifetimeSecs uint) (string, error)
	// DecryptSSOHandle decrypts the relay state from the IDP auth response into
	// an ssoHandle.
	DecryptSSOHandle(encrypted string, key []byte) (string, error)
}

// NewSessionStore creates a SessionStore
func NewSessionStore(pool *redis.Pool) SessionStore {
	return &store{pool}
}

type store struct {
	pool *redis.Pool
}

func (s *store) CreateSession(ssoHandle, originalURL string, lifetimeSecs uint) error {
	if len(ssoHandle) < 8 {
		return errors.New("session handle must be 8 or more characters in length")
	}
	conn := s.pool.Get()
	defer conn.Close()
	// make sure key isn't already in use, this prevents
	// clobbering keys used for something else
	exists, err := redis.Int(conn.Do("EXISTS", ssoHandle))
	if err != nil {
		return err
	}
	if exists == 1 {
		return errors.Errorf("session key '%s' is already in use", ssoHandle)
	}

	expiresAt := time.Now().Add(time.Duration(lifetimeSecs) * time.Second)
	sess := Session{OriginalURL: originalURL, ExpiresAt: expiresAt}
	var writer bytes.Buffer
	err = json.NewEncoder(&writer).Encode(sess)
	if err != nil {
		return err
	}
	_, err = conn.Do("SETEX", ssoHandle, lifetimeSecs, writer.String())
	if err != nil {
		return err
	}
	return err
}

var ErrSessionNotFound = errors.New("session not found")

func (s *store) GetSession(ssoHandle string) (*Session, error) {
	conn := s.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", ssoHandle))
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var sess Session
	reader := bytes.NewBufferString(val)
	err = json.NewDecoder(reader).Decode(&sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *store) UpdateSession(ssoHandle, userName string) (*Session, error) {
	sess, err := s.GetSession(ssoHandle)
	if err != nil {
		return nil, errors.Wrap(err, "updating session")
	}
	conn := s.pool.Get()
	defer conn.Close()
	// SET resets / removes lifetime on object
	// so we have to do some math so we can
	// set remaining lifetime based on the original setting
	lifetime := getRemainingLife(sess.ExpiresAt)
	sess.UserName = userName
	var writer bytes.Buffer
	err = json.NewEncoder(&writer).Encode(sess)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("SETEX", ssoHandle, lifetime, writer.String())
	return sess, err
}

func (s *store) ExpireSession(ssoHandle string) error {
	conn := s.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", ssoHandle)
	return err
}

// we encrypt the ssoHandle because to avoid the possibility that a black hat could intercept it
// and use it to call ssologin and get an auth token
func (s *store) EncryptSSOHandle(ssoHandle string, key []byte, lifetimeSecs uint) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// create a different initialization vector for each relay state and store
	// it in redis
	iv := getIV()
	cleartext := []byte(ssoHandle)
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(cleartext))
	encrypter.XORKeyStream(encrypted, cleartext)
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	escaped := url.QueryEscape(encoded)
	_, err = conn.Do("SETEX", escaped, lifetimeSecs, iv)
	if err != nil {
		return "", err
	}
	return escaped, nil
}

func (s *store) DecryptSSOHandle(relayState string, key []byte) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()
	// get initialization vector
	iv, err := redis.Bytes(conn.Do("GET", relayState))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	unescaped, err := url.QueryUnescape(relayState)
	if err != nil {
		return "", err
	}
	unencoded, err := base64.StdEncoding.DecodeString(unescaped)
	if err != nil {
		return "", err
	}
	encrypted := []byte(unencoded)
	decrypter := cipher.NewCFBDecrypter(block, iv)
	unencrypted := make([]byte, len(encrypted))
	decrypter.XORKeyStream(unencrypted, encrypted)
	return string(unencrypted), nil
}

const (
	ivSize  = 16
	keySize = 32
	charset = "abcdefghijklmnopqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func getIV() []byte {
	iv := make([]byte, ivSize)
	rand.Read(iv)
	for i := range iv {
		iv[i] = charset[int(iv[i])%len(charset)]
	}
	return iv
}

// CreateAESKey creates an AES crypto key of the appropriate size
func CreateAESKey() []byte {
	key := make([]byte, keySize)
	rand.Read(key)
	return key
}

func getRemainingLife(expiry time.Time) int {
	remaining := int(time.Until(expiry).Seconds())
	if remaining <= 0 {
		return 1
	}
	return remaining
}
