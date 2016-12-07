package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_20161118193812, Down_20161118193812)
}

func Up_20161118193812(tx *sql.Tx) error {
	sqlStatement := "CREATE TABLE `app_configs` (" +
		"`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT," +
		"`org_name` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`org_logo_url` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`kolide_server_url` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`smtp_configured` TINYINT(1) NOT NULL DEFAULT FALSE," +
		"`smtp_sender_address` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`smtp_server` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`smtp_port` INT UNSIGNED NOT NULL DEFAULT 465," +
		"`smtp_authentication_type` INT UNSIGNED NOT NULL DEFAULT 0," +
		"`smtp_enable_ssl_tls` TINYINT(1) NOT NULL DEFAULT TRUE," +
		"`smtp_authentication_method` INT UNSIGNED NOT NULL DEFAULT 0," +
		"`smtp_domain` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`smtp_user_name` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`smtp_password` VARCHAR(255) NOT NULL DEFAULT ''," +
		"`smtp_verify_ssl_certs` TINYINT(1) NOT NULL DEFAULT TRUE, " +
		"`smtp_enable_start_tls` TINYINT(1) NOT NULL DEFAULT TRUE, " +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	_, err := tx.Exec(sqlStatement)
	return err

}

func Down_20161118193812(tx *sql.Tx) error {
	sqlStatement := "DROP TABLE IF EXISTS `app_configs`;"
	_, err := tx.Exec(sqlStatement)
	return err
}
