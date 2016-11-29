package inmem

import (
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (orm *Datastore) SessionByKey(key string) (*kolide.Session, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, session := range orm.sessions {
		if session.Key == key {
			return session, nil
		}
	}
	return nil, notFound("Session", 0)
}

func (orm *Datastore) SessionByID(id uint) (*kolide.Session, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if session, ok := orm.sessions[id]; ok {
		return session, nil
	}
	return nil, notFound("Session", id)
}

func (orm *Datastore) ListSessionsForUser(id uint) ([]*kolide.Session, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	var sessions []*kolide.Session
	for _, session := range orm.sessions {
		if session.UserID == id {
			sessions = append(sessions, session)
		}
	}
	if len(sessions) == 0 {
		return nil, notFound("Session", 0)
	}
	return sessions, nil
}

func (orm *Datastore) NewSession(session *kolide.Session) (*kolide.Session, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	session.ID = orm.nextID(session)
	orm.sessions[session.ID] = session
	if err := orm.MarkSessionAccessed(session); err != nil {
		return nil, err
	}

	return session, nil

}

func (orm *Datastore) DestroySession(session *kolide.Session) error {
	if _, ok := orm.sessions[session.ID]; !ok {
		return notFound("Session", session.ID)
	}
	delete(orm.sessions, session.ID)
	return nil
}

func (orm *Datastore) DestroyAllSessionsForUser(id uint) error {
	for _, session := range orm.sessions {
		if session.UserID == id {
			delete(orm.sessions, session.ID)
		}
	}
	return nil
}

func (orm *Datastore) MarkSessionAccessed(session *kolide.Session) error {
	session.AccessedAt = time.Now().UTC()
	if _, ok := orm.sessions[session.ID]; !ok {
		return notFound("Session", session.ID)
	}
	orm.sessions[session.ID] = session
	return nil
}

// TODO test session validation(expiration)
