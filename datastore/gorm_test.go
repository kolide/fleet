package datastore

import (
	"fmt"
	"os"
	"testing"

	"github.com/kolide/kolide-ose/app"
)

// TestUser tests the UserStore interface
// this test uses the MySQL GORM backend
func TestUserMySQLGORM(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	db := setupMySQLGORM(t)
	defer teardownMySQLGORM(t, db)

	testUser(t, db)
}

func setupMySQLGORM(t *testing.T) app.Datastore {
	// TODO use ENV vars from docker config
	user := "kolide"
	password := "kolide"
	host := "127.0.0.1:3306"
	dbName := "kolide"

	conn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err := New("gorm", conn, LimitAttempts(1))
	if err != nil {
		t.Fatal(err)
	}

	backend := db.(gormDB)
	if err := backend.migrate(); err != nil {
		t.Fatal(err)
	}

	return db
}

func teardownMySQLGORM(t *testing.T, db app.Datastore) {
	backend := db.(gormDB)
	if err := backend.rollback(); err != nil {
		t.Fatal(err)
	}
}
