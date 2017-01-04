package tables

import "github.com/kolide/goose"

var (
	Client = goose.Client{
		TableName: "migration_status_tables",
		Dialect:   goose.MySqlDialect{},
	}
)
