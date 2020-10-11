package tables

import (
	"fmt"

	"database/sql"
	"github.com/kolide/goose"
)

func init() {
	goose.AddMigration(Up20201011162341, Down20201011162341)
}

func cleanupSoftDeleteFields(tx *sql.Tx, dbTable string) error {
	deleteStmt := fmt.Sprintf("DELETE FROM `%s` WHERE deleted;", dbTable)

	_, err := tx.Exec(deleteStmt)
	if err != nil {
		return err
	}

	for _, column := range []string{"deleted", "deleted_at"} {
		alterStmt := fmt.Sprint(
			"ALTER TABLE `%s` DROP COLUMN `%s`;", dbTable, column)

		_, err := tx.Exec(alterStmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func addSoftDeleteFields(tx *sql.Tx, dbTable string) error {
	addDeletedStmt := fmt.Sprint(
		"ALTER TABLE `%s` "+
			"ADD COLUMN `deleted` TINYINT(1) NOT NULL DEFAULT FALSE;",
		dbTable)
	_, err := tx.Exec(addDeletedStmt)
	if err != nil {
		return err
	}

	addDeletedAtStmt := fmt.Sprint(
		"ALTER TABLE `%s` "+
			"ADD COLUMN `deleted_at` TIMESTAMP NULL DEFAULT NULL;",
		dbTable)
	_, err = tx.Exec(addDeletedAtStmt)
	if err != nil {
		return err
	}

	return nil
}

func getTablesForCleanupSoftDeletedColumnsMigration() []string {
	tables := []string{
		"distributed_query_campaigns",
		"labels",
		"invites",
		"hosts",
		"packs",
		"queries",
		"scheduled_queries",
		"users",
	}
	return tables
}

func Up20201011162341(tx *sql.Tx) error {
	tables := getTablesForCleanupSoftDeletedColumnsMigration()

	for _, table := range tables {
		err := cleanupSoftDeleteFields(tx, table)
		if err != nil {
			return err
		}
	}

	return nil
}

func Down20201011162341(tx *sql.Tx) error {
	tables := getTablesForCleanupSoftDeletedColumnsMigration()

	for _, table := range tables {
		err := addSoftDeleteFields(tx, table)
		if err != nil {
			return err
		}
	}

	return nil
}
