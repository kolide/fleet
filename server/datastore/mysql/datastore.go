package mysql

import (
	"errors"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // db driver
	"github.com/jmoiron/sqlx"
)

var ErrDatabaseConnection = errors.New("Database connection attempts failed")

// Datastore is an implementation of kolide.Datastore interface backed by
// MySQL
type Datastore struct {
	db     *sqlx.DB
	logger *log.Logger
}

// NewDatastore creates an MySQL datastore.
func NewDatastore(dbConnectString string, opts ...DBOption) (ds *Datastore, e error) {

	var options dbOptions

	for _, setOpt := range opts {
		setOpt(&options)
	}

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

		logger.Printf("Connect attempt %d failed. %s\n", attempt, e.Error())
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
