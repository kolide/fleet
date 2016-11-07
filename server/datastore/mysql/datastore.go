package mysql

import (
	"strings"

	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql" // db driver
	"github.com/jmoiron/sqlx"
)

// Datastore is an implementation of kolide.Datastore interface backed by
// MySQL
type Datastore struct {
	db     *sqlx.DB
	logger log.Logger
}

// New creates an MySQL datastore.
func New(dbConnectString string, opts ...DBOption) (*Datastore, error) {
	var (
		ds  *Datastore
		err error
		db  *sqlx.DB
	)

	options := dbOptions{
		maxAttempts: defaultMaxAttempts,
		logger:      log.NewNopLogger(),
	}

	for _, setOpt := range opts {
		setOpt(&options)
	}

	for attempt := 0; attempt < options.maxAttempts; attempt++ {
		if db, err = sqlx.Connect("mysql", dbConnectString); err == nil {
			break
		}
	}

	if db == nil {
		return nil, err
	}

	ds = &Datastore{db, options.logger}

	ds.log("Datastore created")

	return ds, nil

}

func (d *Datastore) Name() string {
	return "mysql"
}

// Migrate creates database
func (d *Datastore) Migrate() error {
	var (
		err error
		sql []byte
	)

	d.log("Begin database migration")

	if sql, err = Asset("db/up.sql"); err != nil {
		return err
	}

	tx := d.db.MustBegin()

	for _, statement := range strings.SplitAfter(string(sql), ";") {
		if _, err = tx.Exec(statement); err != nil {
			if err.Error() != "Error 1065: Query was empty" {
				tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	d.log("w00t! Migration succeeded")
	return nil

}

// Drop removes database
func (d *Datastore) Drop() error {
	var (
		sql []byte
		err error
	)

	d.log("Dropping database")

	if sql, err = Asset("db/down.sql"); err != nil {
		return err
	}

	tx := d.db.MustBegin()

	for _, statement := range strings.SplitAfter(string(sql), ";") {
		if _, err = tx.Exec(statement); err != nil {
			if err.Error() != "Error 1065: Query was empty" {
				tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	d.log("Database drop succeeds")
	return nil

}

// Close frees resources associated with underlying mysql connection
func (d *Datastore) Close() error {
	return d.db.Close()
}

func (d *Datastore) log(msg string) {
	d.logger.Log("comp", d.Name(), "msg", msg)
}
