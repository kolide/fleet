package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_20170105151732, Down_20170105151732)
}

func Up_20170105151732(tx *sql.Tx) error {
	sqlStatement := "CREATE UNIQUE INDEX idx_query_unique_name " +
		" ON `queries` (`name` ASC);"
	_, err := tx.Exec(sqlStatement)
	return err
}

func Down_20170105151732(tx *sql.Tx) error {
	sqlStatement := "DROP INDEX idx_query_unique_name ON `queries`;"
	_, err := tx.Exec(sqlStatement)
	return err
}
