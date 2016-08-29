package kolide

type OsqueryStore interface {
	OsqueryServerStore
	OsqueryConfigStore
}

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
