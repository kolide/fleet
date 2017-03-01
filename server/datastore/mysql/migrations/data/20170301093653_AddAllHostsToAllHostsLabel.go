package data

import (
	"database/sql"

	"github.com/kolide/kolide/server/kolide"
	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20170301093653, Down_20170301093653)
}

func Up_20170301093653(tx *sql.Tx) error {
	// Get the 'All Hosts' label ID
	var allHostsID uint
	err := tx.QueryRow(`
		 SELECT id FROM labels
                 WHERE name = 'All Hosts'
                 AND label_type = ?
`,
		kolide.LabelTypeBuiltIn).
		Scan(&allHostsID)
	if err != nil {
		return errors.Wrap(err, "finding 'All Hosts' label")
	}

	// Insert any host not currently in that label into the label
	_, err = tx.Exec(`
		 INSERT IGNORE INTO label_query_executions (
                         host_id,
                         label_id,
                         matches
                 ) SELECT id as host_id, ?, true FROM hosts
`,
		allHostsID)
	if err != nil {
		return errors.Wrap(err, "adding hosts to 'All Hosts'")
	}

	return nil
}

func Down_20170301093653(tx *sql.Tx) error {
	// This operation not reversible
	return nil
}
