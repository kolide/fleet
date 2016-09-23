package mock

import "github.com/kolide/kolide-ose/kolide"

var _ kolide.SessionStore = (*SessionStore)(nil)

type SessionStore struct {
	FindSessionByKeyFunc        func(key string) (*kolide.Session, error)
	FindSessionByKeyFuncInvoked bool

	FindSessionByIDFunc        func(id uint) (*kolide.Session, error)
	FindSessionByIDFuncInvoked bool

	FindAllSessionsForUserFunc        func(id uint) ([]*kolide.Session, error)
	FindAllSessionsForUserFuncInvoked bool

	NewSessionFunc        func(session *kolide.Session) (*kolide.Session, error)
	NewSessionFuncInvoked bool

	DestroySessionFunc        func(session *kolide.Session) error
	DestroySessionFuncInvoked bool

	DestroyAllSessionsForUserFunc        func(id uint) error
	DestroyAllSessionsForUserFuncInvoked bool

	MarkSessionAccessedFunc        func(session *kolide.Session) error
	MarkSessionAccessedFuncInvoked bool
}

func (ds *SessionStore) FindSessionByKey(key string) (*kolide.Session, error) {
	ds.FindSessionByKeyFuncInvoked = true
	return ds.FindSessionByKeyFunc(key)
}

func (ds *SessionStore) FindSessionByID(id uint) (*kolide.Session, error) {
	ds.FindSessionByIDFuncInvoked = true
	return ds.FindSessionByIDFunc(id)
}

func (ds *SessionStore) FindAllSessionsForUser(id uint) ([]*kolide.Session, error) {
	ds.FindAllSessionsForUserFuncInvoked = true
	return ds.FindAllSessionsForUserFunc(id)
}

func (ds *SessionStore) NewSession(session *kolide.Session) (*kolide.Session, error) {
	ds.NewSessionFuncInvoked = true
	return ds.NewSessionFunc(session)
}

func (ds *SessionStore) DestroySession(session *kolide.Session) error {
	ds.DestroySessionFuncInvoked = true
	return ds.DestroySessionFunc(session)
}

func (ds *SessionStore) DestroyAllSessionsForUser(id uint) error {
	ds.DestroyAllSessionsForUserFuncInvoked = true
	return ds.DestroyAllSessionsForUserFunc(id)
}

func (ds *SessionStore) MarkSessionAccessed(session *kolide.Session) error {
	ds.MarkSessionAccessedFuncInvoked = true
	return ds.MarkSessionAccessedFunc(session)
}
