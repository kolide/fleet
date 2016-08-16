package datastore

import (
	"fmt"
	"log"
	"testing"
)

// test new conn
func TestNew(t *testing.T) {
	testNewGormConn(t)
}

func testNewGormConn(t *testing.T) {
	connString := map[string]string{
		"failing": fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", "tester", "secret", "kolide.nothere:2022", "kolide"),
	}
	var connTests = []struct {
		conn    string
		willErr bool
	}{
		{connString["failing"], true},
	}

	for _, tt := range connTests {
		_, err := New("gorm", tt.conn, LimitAttempts(1))
		if tt.willErr {
			if err == nil {
				log.Fatal("expected err, but got nil")
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
