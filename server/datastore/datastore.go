// Package datastore implements Kolide's interactions with the database backend
package datastore

import (
	"errors"
	"fmt"

	"github.com/WatchBeam/clock"
	"github.com/kolide/kolide-ose/server/datastore/mysql"
	"github.com/kolide/kolide-ose/server/kolide"
)

var (
	// ErrNotFound is returned when the datastore resource cannot be found
	ErrNotFound = errors.New("resource not found")

	// ErrExists is returned when creating a datastore resource that already exists
	ErrExists = errors.New("resource already created")
)

// New creates a kolide.Datastore with a database connection
// Use DBOption to pass optional arguments
func New(driver, conn string, opts ...mysql.DBOption) (kolide.Datastore, error) {

	switch driver {
	case "mysql":
		ds, err := mysql.New(conn, clock.C, opts...)

		if err != nil {
			return nil, err
		}

		if err = ds.Migrate(); err != nil {
			return nil, err
		}

		return ds, nil

	case "inmem":
		ds := &inmem{
			Driver: "inmem",
		}

		err := ds.Migrate()
		if err != nil {
			return nil, err
		}

		return ds, nil
	default:
		return nil, fmt.Errorf("unsupported datastore driver %s", driver)
	}
}
