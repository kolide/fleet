package mysql

import (
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewScheduledQuery(sq *kolide.PackQuery) (*kolide.PackQuery, error) {
	sql := `
	    INSERT INTO pack_queries (
		    pack_id,
			query_id,
			snapshot,
			differential,
			` + "`interval`" + `,
			platform,
			version,
			shard
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
	result, err := d.db.Exec(sql, sq.PackID, sq.QueryID, sq.Snapshot, sq.Differential, sq.Interval, sq.Platform, sq.Version, sq.Shard)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	id, _ := result.LastInsertId()
	sq.ID = uint(id)
	return sq, nil
}

func (d *Datastore) SaveScheduledQuery(sq *kolide.PackQuery) (*kolide.PackQuery, error) {
	sql := `
		UPDATE pack_queries
			SET pack_id = ?, query_id = ?, interval = ?, snapshot = ?, differential = ?, platform = ?, version = ?, shard = ?
			WHERE id = ? AND NOT deleted
	`
	_, err := d.db.Exec(sql, sq.PackID, sq.QueryID, sq.Interval, sq.Snapshot, sq.Differential, sq.Platform, sq.Version, sq.Shard, sq.ID)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	return sq, nil
}

func (d *Datastore) DeleteScheduledQuery(id uint) error {
	sql := `
		UPDATE pack_queries
			SET deleted_at = ?, deleted = ?
			WHERE id = ?
	`
	_, err := d.db.Exec(sql, d.clock.Now(), true, id)
	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

func (d *Datastore) ScheduledQuery(id uint) (*kolide.PackQuery, error) {
	sql := `
		SELECT * FROM pack_queries WHERE id = ? AND NOT deleted
	`
	sq := &kolide.PackQuery{}
	if err := d.db.Get(sq, sql, id); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return sq, nil
}

func (d *Datastore) ListScheduledQueriesInPack(id uint, opts kolide.ListOptions) ([]*kolide.PackQuery, error) {
	sql := `
		SELECT * FROM pack_queries WHERE pack_id = ? AND NOT deleted
	`
	sql = appendListOptionsToSQL(sql, opts)
	results := []*kolide.PackQuery{}

	if err := d.db.Select(&results, sql); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return results, nil
}
