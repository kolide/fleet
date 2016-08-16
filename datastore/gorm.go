package datastore

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/app"
	"github.com/kolide/kolide-ose/sessions"
)

var tables = [...]interface{}{
	&app.User{},
	&app.PasswordResetRequest{},
	&sessions.Session{},
	&app.ScheduledQuery{},
	&app.Pack{},
	&app.DiscoveryQuery{},
	&app.Host{},
	&app.Label{},
	&app.Option{},
	&app.Decorator{},
	&app.Target{},
	&app.DistributedQuery{},
	&app.Query{},
	&app.DistributedQueryExecution{},
}

type gormDB struct {
	DB *gorm.DB
}

// NewUser creates a new user in the gorm backend
func (orm gormDB) NewUser(user *app.User) (*app.User, error) {
	err := orm.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (orm gormDB) migrate() error {
	var err error
	for _, table := range tables {
		err = orm.DB.AutoMigrate(table).Error
	}
	return err
}

func (orm gormDB) rollback() error {
	var err error
	for _, table := range tables {
		err = orm.DB.DropTableIfExists(table).Error
	}
	return err
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
