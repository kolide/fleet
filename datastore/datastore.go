// Package datastore implements Kolide's interactions with the database backend
package datastore

import "github.com/kolide/kolide-ose/app"

type dbOptions struct {
	maxAttempts int
	db          app.Datastore
}

// DBOption is used to pass optional arguments to a database connection
type DBOption func(o *dbOptions) error

// LimitAttempts sets number of maximum connection attempts
func LimitAttempts(attempts int) DBOption {
	return func(o *dbOptions) error {
		o.maxAttempts = attempts
		return nil
	}
}

// datastore allows you to pass your own datastore
// this option can be used to pass a specific testing implementation
func datastore(db app.Datastore) DBOption {
	return func(o *dbOptions) error {
		o.db = db
		return nil
	}
}

// New creates a Datastore with a database connection
// Use DBOption to pass optional arguments
func New(driver, conn string, opts ...DBOption) (app.Datastore, error) {
	opt := &dbOptions{
		maxAttempts: 15, // default attempts
	}
	for _, option := range opts {
		if err := option(opt); err != nil {
			return nil, err
		}
	}

	// check if datastore is already present
	if opt.db != nil {
		return opt.db, nil
	}

	var db app.Datastore
	switch driver {
	case "gorm":
		db, err := openGORM("mysql", conn, opt.maxAttempts)
		if err != nil {
			return nil, err
		}
		ds := gormDB{DB: db}
		if err := ds.migrate(); err != nil {
			return nil, err
		}
		return ds, nil
	}
	return db, nil
}
