package kolide

// Datastore combines all the interfaces in the Kolide DAL
type Datastore interface {
	UserStore
	OsqueryStore
	EmailStore
	SessionStore
	Name() string
	Drop() error
	Migrate() error
}

// UserStore contains methods for managing users in a datastore
type UserStore interface {
	NewUser(user *User) (*User, error)
	User(username string) (*User, error)
	UserByID(id uint) (*User, error)
	SaveUser(user *User) error
}
