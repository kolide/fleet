package data

import (
	"database/sql"

	"github.com/kolide/fleet/server/kolide"
)

func init() {
	MigrationClient.AddMigration(Up_20170724124700, Down_20170724124700)
}

func Up_20170724124700(tx *sql.Tx) error {
	sqlStatement := `
		INSERT INTO options (
			name,
			type,
			value,
			read_only
		) VALUES (?, ?, ?, ?)
	`

	_, err := tx.Exec(sqlStatement, "file_events", kolide.OptionTypeString, kolide.OptionValue{Val: nil}, kolide.NotReadOnly)
	if err != nil {
		return err
	}

	return nil
}

func Down_20170724124700(tx *sql.Tx) error {
	sqlStatement := `
		DELETE FROM options
		WHERE name = ?
	`
	_, err := tx.Exec(sqlStatement, "file_events")
	if err != nil {
		return err
	}
	return nil
}
