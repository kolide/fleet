package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) NewAppConfig(info *kolide.AppConfig) (*kolide.AppConfig, error) {
	var (
		err    error
		result sql.Result
	)

	err = d.db.Get(info, "SELECT * FROM app_configs LIMIT 1")
	switch err {
	case sql.ErrNoRows:
		result, err = d.db.Exec(
			"INSERT INTO app_configs (org_name, org_logo_url, kolide_server_url) VALUES (?, ?, ?)",
			info.OrgName, info.OrgLogoURL, info.KolideServerURL,
		)
		if err != nil {
			return nil, errors.Wrap(err, "insert new AppConfig")
		}

		info.ID, _ = result.LastInsertId()
		return info, nil
	case nil:
		if err := d.SaveAppConfig(info); err != nil {
			return nil, errors.Wrap(err, "save AppConfig")
		}
		return info, nil
	default:
		return nil, errors.Wrap(err, "get AppConfig")
	}
}

func (d *Datastore) AppConfig() (*kolide.AppConfig, error) {
	var config kolide.AppConfig
	if err := d.db.Get(&config, "SELECT * FROM app_configs LIMIT 1"); err != nil {
		return nil, errors.Wrap(err, "get AppConfig")
	}
	return &config, nil
}

func (d *Datastore) SaveAppConfig(info *kolide.AppConfig) error {
	if _, err := d.db.Exec(
		"UPDATE app_configs SET org_name = ?, org_logo_url = ?, kolide_server_url = ? WHERE id = ?",
		info.OrgName, info.OrgLogoURL, info.KolideServerURL, info.ID,
	); err != nil {
		return errors.Wrap(err, "save AppConfig")
	}
	return nil
}
