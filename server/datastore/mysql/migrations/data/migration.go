package data

import "github.com/kolide/goose"

var (
	Client = goose.Client{
		TableName: "migration_status_data",
		Dialect:   goose.MySqlDialect{},
	}
)
