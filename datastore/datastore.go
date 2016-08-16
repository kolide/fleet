// Package datastore implements Kolide's interactions with the database backend
package datastore

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/app"
)

// Datastore combines all methods for backend interactions
type Datastore interface {
	app.UserStore
}

type gormDB struct {
	DB *gorm.DB
}

// NewUser creates a new user in the gorm backend
func (db gormDB) NewUser(user *app.User) (*app.User, error) {
	panic("not implemented")
}

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
func New(driver, conn string, opts ...DBOption) (Datastore, error) {
	opt := &dbOptions{
		maxAttempts: 15, // default attempts
	}
	for _, option := range opts {
		if err := option(opt); err != nil {
			return nil, err
		}
	}

	switch driver {
	case "gorm":
		db, err := openGORM("mysql", conn, opt.maxAttempts)
		if err != nil {
			return nil, err
		}
		return gormDB{DB: db}, nil
	}
	return nil, nil
}

// create connection with mysql backend, using a backoff timer and maxAttempts
func openGORM(driver, conn string, maxAttempts int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(driver, conn)
		if err == nil {
			break
		} else {
			if err.Error() == "invalid database source" {
				return nil, err
			}
			// TODO: use a logger
			fmt.Printf("could not connect to mysql: %v\n", err)
			time.Sleep(time.Duration(attempts) * time.Second)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql backend, err = %v", err)
	}
	return db, nil
}
