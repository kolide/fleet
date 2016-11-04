package mysql

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func setupMySQL(t *testing.T) (ds *Datastore, teardown func()) {
	var (
		user     = "kolide"
		password = "kolide"
		dbName   = "kolide"
		host     = "127.0.0.1"
	)

	if h, ok := os.LookupEnv("MYSQL_PORT_3306_TCP_ADDR"); ok {
		host = h
	}

	connString := fmt.Sprintf("%s:%s@(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	fmt.Println(connString)

	ds, err := New(connString, Logger(log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))))
	require.Nil(t, err)
	teardown = func() {
		ds.db.Close()
	}

	return ds, teardown
}

func TestDBMigrateAndDrop(t *testing.T) {
	if _, ok := os.LookupEnv("MYSQL_TEST"); !ok {
		t.SkipNow()
	}

	ds, teardown := setupMySQL(t)
	defer teardown()

	err := ds.Migrate()
	require.Nil(t, err)

	err = ds.Drop()
	require.Nil(t, err)

}
