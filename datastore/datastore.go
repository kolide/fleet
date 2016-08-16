// Package datastore implements Kolide's interactions with the database backend
package datastore

import "github.com/kolide/kolide-ose/app"

type dbOptions struct {
	maxAttempts int
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

	var db app.Datastore
	switch driver {
	case "gorm":
		db, err := openGORM("mysql", conn, opt.maxAttempts)
		if err != nil {
			return nil, err
		}
		return gormDB{DB: db}, nil
	}
	return db, nil
}
