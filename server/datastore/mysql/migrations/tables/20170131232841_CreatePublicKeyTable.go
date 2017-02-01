package tables

import (
	"database/sql"
)

func init() {
	MigrationClient.AddMigration(Up_20170131232841, Down_20170131232841)
}

func Up_20170131232841(tx *sql.Tx) error {
	create := "CREATE TABLE `public_keys` ( " +
		"`hash` char(64) NOT NULL DEFAULT '', " +
		"`key` text NOT NULL, " +
		"PRIMARY KEY (`hash`) " +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	_, err := tx.Exec(create)
	if err != nil {
		return err
	}
	return nil
}

func Down_20170131232841(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS `public_keys`;")
	if err != nil {
		return err
	}
	return nil
}
