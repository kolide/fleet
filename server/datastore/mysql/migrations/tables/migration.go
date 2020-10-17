package tables

import "github.com/kolide/goose"

var (
	MigrationClient = goose.New("migration_status_tables", goose.MySqlDialect{})
)

const (
	// Unknown column '%s' in '%s'
	//
	// This error occurs when you try to use a column that doesn't exist in the
	// WHERE clause of a statement.
	ER_BAD_FIELD_ERROR = 1054

	// Message: Can't DROP '%s'; check that column/key exists
	//
	// This error occurs when you try to use a column that doesn't exist in the
	// DROP clause of an ALTER statement.
	ER_CANT_DROP_FIELD_OR_KEY = 1091
)
