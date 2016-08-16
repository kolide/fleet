package datastore

import (
	"fmt"
	"testing"

	"github.com/kolide/kolide-ose/app"
)

func TestNewUser(t *testing.T) {
	db := newDB(t)
	defer teardown(t, db)

	var userTests = []struct {
		username, password, email string
		isAdmin, passwordReset    bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "foobar", "jason@kolide.co", true, false},
	}

	for _, tt := range userTests {
		u, err := app.NewUser(tt.username, tt.password, tt.email, tt.isAdmin, tt.passwordReset)
		if err != nil {
			t.Fatal(err)
		}

		user, err := db.NewUser(u)
		if err != nil {
			t.Fatal(err)
		}

		if user.Username != tt.username {
			t.Errorf("expected username: %s, got %s", user.Username, tt.username)
		}

		if user.Email != tt.email {
			t.Errorf("expected email: %s, got %s", user.Email, tt.email)
		}
		if !user.Admin {
			t.Errorf("expected email: %s, got %s", user.Email, tt.email)
		}

		backend := db.(gormDB)
		var verify app.User
		backend.DB.Where("username = ?", tt.username).First(&verify)
		if verify.ID != user.ID {
			t.Fatal("Couldn't select user back from database")
		}
	}

}

// new db creates a test datastore
func newDB(t *testing.T) app.Datastore {
	user := "kolide"
	password := "kolide"
	host := "127.0.0.1:3306"
	dbName := "kolide"

	conn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err := openGORM("mysql", conn, 1)
	if err != nil {
		fmt.Println(conn)
		t.Fatal(err)
	}
	ds := gormDB{DB: db}
	if err := ds.migrate(); err != nil {
		t.Fatal(err)
	}
	return ds
}

func teardown(t *testing.T, ds app.Datastore) {
	backend := ds.(gormDB)
	if err := backend.rollback(); err != nil {
		t.Fatal(err)
	}
}
