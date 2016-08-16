package datastore

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/app"
)

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
