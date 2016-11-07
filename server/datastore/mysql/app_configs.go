package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewOrgInfo(info *kolide.OrgInfo) (*kolide.OrgInfo, error) {
	var (
		err    error
		result sql.Result
	)

	err = d.db.Get(info, "SELECT * FROM org_infos WHERE org_name = ? LIMIT 1", info.OrgName)
	switch err {
	case sql.ErrNoRows:
		result, err = d.db.Exec(
			"INSERT INTO org_infos (org_name, org_logo_url) VALUES (?, ?)",
			info.OrgName, info.OrgLogoURL,
		)

		if err != nil {
			return nil, err
		}
		info.ID, _ = result.LastInsertId()
		return info, nil
	case nil:
		return info, d.SaveOrgInfo(info)
	default:
		return nil, err
	}
}

func (d *Datastore) OrgInfo() (*kolide.OrgInfo, error) {
	info := &kolide.OrgInfo{}
	err := d.db.Get(info, "SELECT * FROM org_infos LIMIT 1")
	return info, err
}

func (d *Datastore) SaveOrgInfo(info *kolide.OrgInfo) error {
	_, err := d.db.Exec(
		"UPDATE org_infos SET org_name = ?, org_logo_url = ? WHERE id = ?",
		info.OrgName, info.OrgLogoURL, info.ID,
	)
	return err
}
