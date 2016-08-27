// Package datastore implements Kolide's interactions with the database backend
package datastore

import (
	"github.com/kolide/kolide-ose/errors"
	"github.com/kolide/kolide-ose/kolide"
)

// New creates a Datastore with a database connection
// Use DBOption to pass optional arguments
func New(driver, conn string, opts ...DBOption) (kolide.Datastore, error) {
	opt := &dbOptions{
		// configure defaults
		maxAttempts:     defaultMaxAttempts,
		sessionLifespan: defaultSessionLifespan,
		sessionKeySize:  defaultSessionKeySize,
	}
	for _, option := range opts {
		if err := option(opt); err != nil {
			return nil, errors.DatabaseError(err)
		}
	}

	// check if datastore is already present
	if opt.db != nil {
		return opt.db, nil
	}

	switch driver {
	case "gorm-mysql":
		db, err := openGORM("mysql", conn, opt.maxAttempts)
		if err != nil {
			return nil, errors.DatabaseError(err)
		}
		ds := gormDB{
			DB:              db,
			Driver:          "mysql",
			sessionKeySize:  opt.sessionKeySize,
			sessionLifespan: opt.sessionLifespan,
		}
		// configure logger
		if opt.logger != nil {
			db.SetLogger(opt.logger)
			db.LogMode(opt.debug)
		}
		if err := ds.Migrate(); err != nil {
			return nil, errors.DatabaseError(err)
		}
		return ds, nil
	case "gorm-sqlite3":
		db, err := openGORM("sqlite3", conn, opt.maxAttempts)
		if err != nil {
			return nil, errors.DatabaseError(err)
		}
		ds := gormDB{
			DB:              db,
			Driver:          "sqlite3",
			sessionKeySize:  opt.sessionKeySize,
			sessionLifespan: opt.sessionLifespan,
		}
		// configure logger
		if opt.logger != nil {
			db.SetLogger(opt.logger)
			db.LogMode(opt.debug)
		}
		if err := ds.Migrate(); err != nil {
			return nil, errors.DatabaseError(err)
		}
		return ds, nil
	case "mock":
		ds := &mockDB{
			Driver:          "mock",
			sessionKeySize:  opt.sessionKeySize,
			sessionLifespan: opt.sessionLifespan,
			users:           make(map[uint]*kolide.User),
		}
		return ds, nil
	default:
		return nil, errors.New("unsupported datastore driver %s", driver)
	}
}
