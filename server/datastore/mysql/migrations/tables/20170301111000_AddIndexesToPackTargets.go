package tables

import "database/sql"

func init() {
	MigrationClient.AddMigration(Up_20170301111000, Down_20170301111000)
}

func Up_20170301111000(tx *sql.Tx) error {
	indexes := []string{
		"ALTER TABLE `pack_targets` ADD INDEX `idx_pack_targets_pack_id` (`pack_id`);",
		"ALTER TABLE `pack_targets` ADD INDEX `idx_pack_targets_target_id` (`target_id`);",
		"ALTER TABLE `pack_targets` ADD INDEX `idx_pack_targets_type` (`type`);",
	}
	for _, q := range indexes {
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func Down_20170301111000(tx *sql.Tx) error {
	indexes := []string{
		"ALTER TABLE `pack_targets` DROP INDEX `idx_pack_targets_pack_id`;",
		"ALTER TABLE `pack_targets` DROP INDEX `idx_pack_targets_target_id`;",
		"ALTER TABLE `pack_targets` DROP INDEX `idx_pack_targets_type`;",
	}
	for _, q := range indexes {
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}
