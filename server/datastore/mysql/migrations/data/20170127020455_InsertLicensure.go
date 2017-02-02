package data

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	"github.com/kolide/kolide/server/datastore/internal/appstate"
)

func init() {
	MigrationClient.AddMigration(Up_20170127020455, Down_20170127020455)
}

// hexadecimal md5 hash grouped by 2 characters separated by colons
func fingerprintMD5(pem string) string {
	var fingerPrint string
	hash := md5.Sum([]byte(pem))
	for _, h := range hash {
		if len(fingerPrint) > 0 {
			fingerPrint += ":"
		}
		fingerPrint += fmt.Sprintf("%.2x", h)
	}
	return fingerPrint
}

func Up_20170127020455(tx *sql.Tx) error {
	for _, pem := range appstate.PublicKeys() {
		fingerPrint := fingerprintMD5(pem)
		_, err := tx.Exec("INSERT INTO public_keys (hash, `key`) VALUES(?, ?);", fingerPrint, pem)
		if err != nil {
			return err
		}
	}
	_, err := tx.Exec(
		"INSERT INTO licenses ( " +
			"id, " +
			"revoked, " +
			"`key`  " +
			") VALUES (1, FALSE, '');",
	)
	return err
}

func Down_20170127020455(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM public_keys;`)
	return err
}
