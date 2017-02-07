package tables

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20170207103441, Down_20170207103441)
}

func Up_20170207103441(tx *sql.Tx) error {
	for _, query := range []string{
		"UPDATE `scheduled_queries` SET `snapshot` = false WHERE `snapshot` is NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `snapshot` BOOL NOT NULL DEFAULT false;",
		"UPDATE `scheduled_queries` SET `removed` = true WHERE `removed` is NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `removed` BOOL NOT NULL DEFAULT true;",
		"UPDATE `scheduled_queries` SET `platform` = '' WHERE `platform` is NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `platform` VARCHAR(255) NOT NULL DEFAULT '';",
		"UPDATE `scheduled_queries` SET `version` = '' WHERE `version` is NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `version` VARCHAR(255) NOT NULL DEFAULT '';",
		"UPDATE `scheduled_queries` SET `shard` = 100 WHERE `shard` is NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `shard` INT(10) UNSIGNED NOT NULL DEFAULT 100;",
	} {
		_, err := tx.Exec(query)
		if err != nil {
			return errors.Wrap(err, query)
		}
	}
	return nil
}

func Down_20170207103441(tx *sql.Tx) error {
	for _, query := range []string{
		"ALTER TABLE `scheduled_queries` MODIFY `snapshot` BOOL DEFAULT NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `removed` BOOL DEFAULT NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `platform` VARCHAR(255) DEFAULT NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `version` VARCHAR(255) DEFAULT NULL;",
		"ALTER TABLE `scheduled_queries` MODIFY `shard` INT(10) UNSIGNED DEFAULT NULL;",
	} {
		_, err := tx.Exec(query)
		if err != nil {
			return errors.Wrap(err, query)
		}
	}
	return nil
}
