package data

import (
	"database/sql"

	"github.com/kolide/goose"
	"github.com/kolide/kolide-ose/server/datastore/internal/appstate"
)

func init() {
	goose.AddMigration(Up_20170127020455, Down_20170127020455)
}

func Up_20170127020455(tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO licensure (id, public_key) VALUES(1, ?)`, appstate.PublicKey)
	return nil
}

func Down_20170127020455(tx *sql.Tx) error {
	return nil
}
