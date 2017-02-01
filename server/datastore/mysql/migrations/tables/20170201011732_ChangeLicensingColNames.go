package tables

import (
	"database/sql"
)

func init() {
	MigrationClient.AddMigration(Up_20170201011732, Down_20170201011732)
}

func Up_20170201011732(tx *sql.Tx) error {
	_, err := tx.Exec(
		"ALTER TABLE `kolide`.`licensure` CHANGE COLUMN `public_key` " +
			"`key` TEXT CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL " +
			" COMMENT '' AFTER `revoked`;",
	)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"ALTER TABLE `kolide`.`licensure` CHANGE COLUMN `license` " +
			"`token` TEXT CHARACTER SET utf8 COLLATE utf8_general_ci NULL " +
			" COMMENT '' AFTER `key`;",
	)
	if err != nil {
		return err
	}
	return nil
}

func Down_20170201011732(tx *sql.Tx) error {
	_, err := tx.Exec(
		"ALTER TABLE `kolide`.`licensure` CHANGE COLUMN `key` " +
			"`public_key` TEXT CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL " +
			" COMMENT '' AFTER `revoked`;",
	)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"ALTER TABLE `kolide`.`licensure` CHANGE COLUMN `token` " +
			"`license` TEXT CHARACTER SET utf8 COLLATE utf8_general_ci NULL " +
			" COMMENT '' AFTER `public_key`;",
	)
	if err != nil {
		return err
	}
	return nil
}
