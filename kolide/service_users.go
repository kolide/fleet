package kolide

// UserService has methods for working with users
type UserService interface {
	UserStore
	SetPassword(userID uint, password string) error
}

func (svc service) NewUser(user *User) (*User, error) {
	err := user.setPassword(string(user.Password), svc.saltKeySize, svc.bcryptCost)
	if err != nil {
		return nil, err
	}

	user, err = svc.db.NewUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc service) User(username string) (*User, error) {
	return svc.db.User(username)
}

func (svc service) UserByID(id uint) (*User, error) {
	return svc.db.UserByID(id)
}

func (svc service) SaveUser(user *User) error {
	return svc.db.SaveUser(user)
}

func (svc service) SetPassword(userID uint, password string) error {
	user, err := svc.UserByID(userID)
	if err != nil {
		return err
	}

	err = user.setPassword(password, svc.saltKeySize, svc.bcryptCost)
	if err != nil {
		return err
	}
	err = svc.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}
