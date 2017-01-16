package demo

import "github.com/kolide/goose"

var (
	MigrationClient = goose.New("migration_status_demo", goose.MySqlDialect{})
)
