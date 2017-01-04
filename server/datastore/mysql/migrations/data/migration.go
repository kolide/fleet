package data

import "github.com/kolide/goose"

var (
	Client = goose.New("migration_status_data", goose.MySqlDialect{})
)
