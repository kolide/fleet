package data

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/datastore/internal/appstate"
)

func init() {
	MigrationClient.AddMigration(Up_20170127020455, Down_20170127020455)
}

func Up_20170127020455(tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO licensure (id, public_key) VALUES(1, ?);`, appstate.PublicKey)
	if err != nil {
		return err
	}
	return nil
}

func Down_20170127020455(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM licensure;`)
	if err != nil {
		return err
	}
	return nil
}
