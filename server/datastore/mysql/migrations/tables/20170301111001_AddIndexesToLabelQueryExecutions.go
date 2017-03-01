package tables

import "database/sql"

func init() {
	MigrationClient.AddMigration(Up_20170301111001, Down_20170301111001)
}

func Up_20170301111001(tx *sql.Tx) error {
	indexes := []string{
		"ALTER TABLE `label_query_executions` ADD INDEX `idx_label_query_executions_host_id` (`host_id`);",
		"ALTER TABLE `label_query_executions` ADD INDEX `idx_label_query_executions_label_id` (`label_id`);",
		"ALTER TABLE `label_query_executions` ADD INDEX `idx_label_query_executions_matches` (`matches`);",
	}
	for _, q := range indexes {
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func Down_20170301111001(tx *sql.Tx) error {
	indexes := []string{
		"ALTER TABLE `label_query_executions` DROP INDEX `idx_label_query_executions_host_id`;",
		"ALTER TABLE `label_query_executions` DROP INDEX `idx_label_query_executions_label_id`;",
		"ALTER TABLE `label_query_executions` DROP INDEX `idx_label_query_executions_matches`;",
	}
	for _, q := range indexes {
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}
