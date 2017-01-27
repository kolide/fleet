package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (ds *Datastore) SaveLicense(jwt string) error {
	sqlStatement := `
    UPDATE licensure SET
      license = ?
    WHERE id = 1
    `
	_, err := ds.db.Exec(sqlStatement, jwt)
	if err != nil {
		return errors.Wrap(err, "saving license")
	}
	return nil
}

func (ds *Datastore) License() (*kolide.License, error) {
	query := `
  SELECT * FROM licensure
    WHERE id = 1
  `
	var license kolide.License
	err := ds.db.Get(&license, query)
	if err != nil {
		return nil, errors.Wrap(err, "fetching license information")
	}
	return &license, nil
}
