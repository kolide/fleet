package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (ds *Datastore) SaveLicense(jwt string) (*kolide.License, error) {
	sqlStatement := `
    UPDATE licensure SET
      license = ?
    WHERE id = 1
    `
	_, err := ds.db.Exec(sqlStatement, jwt)
	if err != nil {
		return nil, errors.Wrap(err, "saving license")
	}
	result, err := ds.License()
	if err != nil {
		return nil, errors.Wrap(err, "fetching license")
	}
	return result, nil
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
	query = `
    SELECT count(*)
      FROM hosts
      WHERE NOT deleted
  `
	err = ds.db.Get(&license.HostCount, query)
	if err != nil {
		return nil, errors.Wrap(err, "fetching host count for license")
	}
	return &license, nil
}
