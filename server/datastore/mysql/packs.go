package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) PackByName(name string) (*kolide.Pack, bool, error) {
	sqlStatement := `
		SELECT *
			FROM packs
			WHERE name = ? AND NOT deleted
	`
	var pack kolide.Pack
	err := d.db.Get(&pack, sqlStatement, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, errors.Wrap(err, "fetching packs by name")
	}

	return &pack, true, nil
}

// NewPack creates a new Pack
func (d *Datastore) NewPack(pack *kolide.Pack) (*kolide.Pack, error) {

	sql := `
		INSERT INTO packs ( name, description, platform, created_by, disabled )
			VALUES ( ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(sql, pack.Name, pack.Description, pack.Platform, pack.CreatedBy, pack.Disabled)
	if err != nil {
		return nil, errors.Wrap(err, "creating new pack")
	}

	id, _ := result.LastInsertId()
	pack.ID = uint(id)
	return pack, nil
}

// SavePack stores changes to pack
func (d *Datastore) SavePack(pack *kolide.Pack) error {

	sql := `
			UPDATE packs
			SET name = ?, platform = ?, disabled = ?, description = ?
			WHERE id = ? AND NOT deleted
	`

	_, err := d.db.Exec(sql, pack.Name, pack.Platform, pack.Disabled, pack.Description, pack.ID)
	if err != nil {
		return errors.Wrap(err, "saving pack with id")
	}

	return nil
}

// DeletePack soft deletes a kolide.Pack so that it won't show up in results
func (d *Datastore) DeletePack(pid uint) error {
	return d.deleteEntity("packs", pid)
}

// Pack fetch kolide.Pack with matching ID
func (d *Datastore) Pack(pid uint) (*kolide.Pack, error) {
	sqlStatement := `SELECT * FROM packs WHERE id = ? AND NOT deleted`
	pack := &kolide.Pack{}
	if err := d.db.Get(pack, sqlStatement, pid); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Pack").WithID(pid)
		}
		return nil, errors.Wrap(err, "fetching pack")
	}

	return pack, nil
}

// ListPacks returns all kolide.Pack records limited and sorted by kolide.ListOptions
func (d *Datastore) ListPacks(opt kolide.ListOptions) ([]*kolide.Pack, error) {
	sqlStatement := `SELECT * FROM packs WHERE NOT deleted`
	sqlStatement = appendListOptionsToSQL(sqlStatement, opt)
	packs := []*kolide.Pack{}
	if err := d.db.Select(&packs, sqlStatement); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Pack")
		}
		return nil, errors.Wrap(err, "error selecting packs")
	}
	return packs, nil
}

// AddLabelToPack associates a kolide.Label with a kolide.Pack
func (d *Datastore) AddLabelToPack(lid uint, pid uint) error {
	sql := `
		INSERT INTO pack_targets ( pack_id,	type, target_id )
			VALUES ( ?, ?, ? )
			ON DUPLICATE KEY UPDATE id=id
	`
	_, err := d.db.Exec(sql, pid, kolide.TargetLabel, lid)
	if err != nil {
		return errors.Wrap(err, "adding label to pack")
	}

	return nil
}

// AddHostToPack associates a kolide.Host with a kolide.Pack
func (d *Datastore) AddHostToPack(hid, pid uint) error {
	sql := `
		INSERT INTO pack_targets ( pack_id, type, target_id )
			VALUES ( ?, ?, ? )
			ON DUPLICATE KEY UPDATE id=id
	`
	_, err := d.db.Exec(sql, pid, kolide.TargetHost, hid)
	if err != nil {
		return errors.Wrap(err, "adding host to pack")
	}

	return nil
}

// ListLabelsForPack will return a list of kolide.Label records associated with kolide.Pack
func (d *Datastore) ListLabelsForPack(pid uint) ([]*kolide.Label, error) {
	sql := `
	SELECT
		l.id,
		l.created_at,
		l.updated_at,
		l.name
	FROM
		labels l
	JOIN
		pack_targets pt
	ON
		pt.target_id = l.id
	WHERE
		pt.type = ?
			AND
		pt.pack_id = ?
	AND NOT l.deleted
	`

	labels := []*kolide.Label{}

	if err := d.db.Select(&labels, sql, kolide.TargetLabel, pid); err != nil {
		return nil, errors.Wrap(err, "finding labels for pack")
	}

	return labels, nil
}

// RemoreLabelFromPack will remove the association between a kolide.Label and
// a kolide.Pack
func (d *Datastore) RemoveLabelFromPack(lid, pid uint) error {
	sql := `
		DELETE FROM pack_targets
			WHERE target_id = ? AND pack_id = ? AND type = ?
	`
	if _, err := d.db.Exec(sql, lid, pid, kolide.TargetLabel); err != nil {
		return errors.Wrap(err, "deleting pack")
	}

	return nil
}

// RemoveHostFromPack will remove the association between a kolide.Host and a
// kolide.Pack
func (d *Datastore) RemoveHostFromPack(hid, pid uint) error {
	sql := `
		DELETE FROM pack_targets
			WHERE target_id = ? AND pack_id = ? AND type = ?
	`
	if _, err := d.db.Exec(sql, hid, pid, kolide.TargetHost); err != nil {
		return errors.Wrap(err, "removing host from pack")
	}

	return nil

}

func (d *Datastore) ListHostsInPack(pid uint, opt kolide.ListOptions) ([]*kolide.Host, error) {
	sqlStatement := `
		SELECT DISTINCT h.*
		FROM hosts h
		JOIN pack_targets pt
		JOIN label_query_executions lqe
		ON (
		  pt.target_id = lqe.label_id
		  AND lqe.host_id = h.id
		  AND lqe.matches
		  AND pt.type = ?
		) OR (
		  pt.target_id = h.id
		  AND pt.type = ?
		)
		WHERE pt.pack_id = ?
	`
	sqlStatement = appendListOptionsToSQL(sqlStatement, opt)
	hosts := []*kolide.Host{}
	if err := d.db.Select(&hosts, sqlStatement, kolide.TargetLabel, kolide.TargetHost, pid); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Hosts")
		}
		return nil, errors.Wrap(err, "list hosts in pack ")
	}
	return hosts, nil
}
