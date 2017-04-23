package tables

import (
	"database/sql"
)

func init() {
	MigrationClient.AddMigration(Up_20170422151643, Down_20170422151643)
}

func Up_20170422151643(tx *sql.Tx) error {
	query := "ALTER TABLE `kolide`.`app_configs` " +
		"ADD COLUMN `aes_key` VARBINARY(128) " +
		"AFTER `osquery_enroll_secret`;"
	_, err := tx.Exec(query)
	return err
}

func Down_20170422151643(tx *sql.Tx) error {
	query := "ALTER TABLE `kolide`.`app_configs` DROP COLUMN `aes_key` ;"
	_, err := tx.Exec(query)
	return err
}
