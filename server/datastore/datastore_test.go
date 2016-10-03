package datastore

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setup creates a datastore for testing
func setup(t *testing.T) kolide.Datastore {
	db, err := gorm.Open("sqlite3", ":memory:")
	require.Nil(t, err)

	ds := gormDB{DB: db, Driver: "sqlite3"}

	err = ds.Migrate()
	assert.Nil(t, err)
	return ds
}

func teardown(t *testing.T, ds kolide.Datastore) {
	err := ds.Drop()
	assert.Nil(t, err)
}
