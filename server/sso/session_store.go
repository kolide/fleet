package sso

import (
	"bytes"

	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

type Session struct {
	OriginalURL string `json:"original_url"`
	UserName    string `json:"user_name"`
}

type SessionStore interface {
	CreateSession(ssoHandle, originalURL string, lifetime int) error
	GetSession(ssoHandle string) (*Session, error)
	UpdateSession(ssoHandle string, session *Session, lifetime int) error
	ExpireSession(ssoHandle string) error
}

func NewSessionStorage(pool *redis.Pool) SessionStore {
	return &store{pool}
}

type store struct {
	pool *redis.Pool
}

func (s *store) CreateSession(ssoHandle, originalURL string, lifetime int) error {
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

	sess := Session{OriginalURL: originalURL}
	var writer bytes.Buffer
	err = json.NewEncoder(&writer).Encode(sess)
	if err != nil {
		return err
	}
	_, err = conn.Do("SETEX", ssoHandle, lifetime, writer.String())
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

func (s *store) UpdateSession(ssoHandle string, session *Session, lifetime int) error {
	conn := s.pool.Get()
	defer conn.Close()
	var writer bytes.Buffer
	err := json.NewEncoder(&writer).Encode(session)
	if err != nil {
		return err
	}
	_, err = conn.Do("SETEX", ssoHandle, lifetime, writer.String())
	return err
}

func (s *store) ExpireSession(ssoHandle string) error {
	conn := s.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", ssoHandle)
	return err

}
