package migration

import (
	"database/sql"
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_20161123011056, Down_20161123011056)
}

func Up_20161123011056(tx *sql.Tx) error {
	// create a column to represent built in types, default is mutable
	_, err := tx.Exec(fmt.Sprintf("ALTER TABLE labels ADD COLUMN label_type INT UNSIGNED NOT NULL DEFAULT %d;", kolide.LabelTypeMutable))

	return err
}

func Down_20161123011056(tx *sql.Tx) error {
	if _, err := tx.Exec("ALTER TABLE labels DROP COLUMN label_type;"); err != nil {
		return err
	}

	return nil
}
