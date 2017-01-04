package tables

import "github.com/kolide/goose"

var (
	Client = goose.New("migration_status_tables", goose.MySqlDialect{})
)
