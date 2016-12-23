package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_20161118212604, Down_20161118212604)
}

func Up_20161118212604(tx *sql.Tx) error {

	_, err := tx.Exec(
		"CREATE TABLE `options` (" +
			"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT," +
			"`name` VARCHAR(255) NOT NULL," +
			"`type` INT UNSIGNED NOT NULL," +
			"`value` VARCHAR(255) DEFAULT NULL," +
			"`read_only` TINYINT(1) NOT NULL DEFAULT FALSE," +
			"PRIMARY KEY (`id`)," +
			"UNIQUE KEY `idx_option_unique_name` (`name`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;",
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"CREATE TABLE `option_values` (" +
			"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
			"`created_at` timestamp DEFAULT CURRENT_TIMESTAMP," +
			"`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"`option_id` INT UNSIGNED NULL," +
			"`value` varchar(255) NOT NULL," +
			"PRIMARY KEY (`id`)," +
			"CONSTRAINT FOREIGN KEY `idx_options_fkey` (`option_id`) " +
			"REFERENCES options(id) " +
>>>>>>> 4b70231267f7ef192c1f7613fb90b85d613315fb
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;",
	)

	return err
}

func Down_20161118212604(tx *sql.Tx) error {
<<<<<<< HEAD

	_, err := tx.Exec("DROP TABLE IF EXISTS `options`;")
=======
	_, err := tx.Exec("DROP TABLE IF EXISTS `option_values`;")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS `options`;")
>>>>>>> 4b70231267f7ef192c1f7613fb90b85d613315fb

	return err
}
