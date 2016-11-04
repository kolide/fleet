package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // db driver
	"github.com/jmoiron/sqlx"
)

const defaultMaxDBConnectionAttempts = 15

var ErrDatabaseConnection = errors.New("Database connection attempts failed")

// Datastore is an implementation of kolide.Datastore interface backed by
// MySQL
type Datastore struct {
	db     *sqlx.DB
	logger *log.Logger
}

// NewDatastore creates an MySQL datastore.
func NewDatastore(dbConnectString string, opts ...DBOption) (ds *Datastore, e error) {
	fmt.Println("Called new ds")
	options := dbOptions{
		maxAttempts: defaultMaxDBConnectionAttempts,
	}

	for _, setOpt := range opts {
		setOpt(&options)
	}

	fmt.Printf("Options %#v\n", options)
	var logger *log.Logger

	if options.logger == nil {
		logger = log.New(os.Stdout, "db", log.LstdFlags)
	} else {
		logger = options.logger
	}

	var db *sqlx.DB

	for attempt := 0; attempt < options.maxAttempts; attempt++ {
		if db, e = sqlx.Connect("mysql", dbConnectString); e == nil {
			break
		}

		fmt.Printf("Connect attempt %d failed. %s\n", attempt, e.Error())
	}

	if db == nil {
		e = ErrDatabaseConnection
	} else {
		ds = &Datastore{db, logger}
	}

	return

}

func (d *Datastore) Name() string {
	return "mysql"
}

// Migrate creates database
func (d *Datastore) Migrate() (e error) {
	d.logger.Println("Begin database migration")

	var sql []byte
	if sql, e = Asset("db/up.sql"); e != nil {
		return
	}

	tx := d.db.MustBegin()

	for _, statement := range strings.SplitAfter(string(sql), ";") {
		if _, e = tx.Exec(statement); e != nil {
			if e.Error() != "Error 1065: Query was empty" {
				tx.Rollback()
				return
			}
		}
	}

	if e = tx.Commit(); e != nil {
		return
	}

	d.logger.Println("w00t! Migration succeeded")
	return

}

// Drop removes database
func (d *Datastore) Drop() (e error) {
	d.logger.Println("Dropping database")

	var sql []byte
	if sql, e = Asset("db/down.sql"); e != nil {
		return
	}

	tx := d.db.MustBegin()

	for _, statement := range strings.SplitAfter(string(sql), ";") {
		if _, e = tx.Exec(statement); e != nil {
			if e.Error() != "Error 1065: Query was empty" {
				tx.Rollback()
				return
			}
		}
	}

	if e = tx.Commit(); e != nil {
		return
	}

	d.logger.Println("Database drop succeeds")
	return

}
