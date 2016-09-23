package mock

import "github.com/kolide/kolide-ose/kolide"

var _ kolide.UserStore = (*UserStore)(nil)

type UserStore struct {
	NewUserFunc        func(user *kolide.User) (*kolide.User, error)
	NewUserFuncInvoked bool

	UserFunc        func(username string) (*kolide.User, error)
	UserFuncInvoked bool

	UsersFunc        func() ([]*kolide.User, error)
	UsersFuncInvoked bool

	UserByEmailFunc        func(email string) (*kolide.User, error)
	UserByEmailFuncInvoked bool

	UserByIDFunc        func(id uint) (*kolide.User, error)
	UserByIDFuncInvoked bool

	SaveUserFunc        func(user *kolide.User) error
	SaveUserFuncInvoked bool
}

func (ds *UserStore) NewUser(user *kolide.User) (*kolide.User, error) {
	ds.NewUserFuncInvoked = true
	return ds.NewUserFunc(user)
}

func (ds *UserStore) User(username string) (*kolide.User, error) {
	ds.UserFuncInvoked = true
	return ds.UserFunc(username)
}

func (ds *UserStore) Users() ([]*kolide.User, error) {
	ds.UsersFuncInvoked = true
	return ds.UsersFunc()
}

func (ds *UserStore) UserByEmail(email string) (*kolide.User, error) {
	ds.UserByEmailFuncInvoked = true
	return ds.UserByEmailFunc(email)
}

func (ds *UserStore) UserByID(id uint) (*kolide.User, error) {
	ds.UserByIDFuncInvoked = true
	return ds.UserByIDFunc(id)
}

func (ds *UserStore) SaveUser(user *kolide.User) error {
	ds.SaveUserFuncInvoked = true
	return ds.SaveUserFunc(user)
}
