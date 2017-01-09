package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_20170108191242, Down_20170108191242)
}

func Up_20170108191242(tx *sql.Tx) error {
	_, err := tx.Exec(
		"CREATE TABLE `decorators` ( " +
			"`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT, " +
			"`created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
			"`updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, " +
			"`query` VARCHAR(255) NOT NULL, " +
			"`type` INT UNSIGNED NOT NULL, " +
			"`interval` INT UNSIGNED NOT NULL, " +
			"PRIMARY KEY (`id`) " +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;",
	)
	return err
}

func Down_20170108191242(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS decorators;")
	return err
}
