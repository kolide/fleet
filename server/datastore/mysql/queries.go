package mysql

import "github.com/kolide/kolide-ose/server/kolide"

// NewQuery creates a Query
func (d *Datastore) NewQuery(query *kolide.Query) (*kolide.Query, error) {
	query.MarkAsCreated(d.clock.Now())
	sql := `
		INSERT INTO queries (created_at, updated_at, name, description, query,
			snapshot, differential, platform, version, ` + "`interval`" + `)
		VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )
	`

	result, err := d.db.Exec(sql, query.CreatedAt, query.UpdatedAt,
		query.Name, query.Description, query.Query, query.Snapshot,
		query.Differential, query.Platform, query.Version, query.Interval)

	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	query.ID = uint(id)
	return query, nil
}

// SaveQuery saves changes to a Query.
func (d *Datastore) SaveQuery(q *kolide.Query) error {
	q.MarkAsUpdated(d.clock.Now())
	// TODO it might be better to use a table alias here to deal with interval
	sql := `
		UPDATE queries
			SET updated_at = ?, name = ?, description = ?, query = ?, ` + "`interval`" + `= ? snapshot = ?,
			 	differential = ?, platform = ?, version = ?
			WHERE id = ? AND NOT deleted
	`
	_, err := d.db.Exec(sql, q.UpdatedAt, q.Name, q.Description, q.Query, q.Interval,
		q.Snapshot, q.Differential, q.Platform, q.Version, q.ID)

	return err
}

// DeleteQuery soft deletes Query identified by Query.ID
func (d *Datastore) DeleteQuery(query *kolide.Query) error {

	query.MarkDeleted(d.clock.Now())
	sql := `
		UPDATE queries
			SET deleted_at = ?, deleted = ?
			WHERE id = ?
	`
	_, err := d.db.Exec(sql, query.DeletedAt, true, query.ID)
	return err
}

// Query returns a single Query identified by id, if such
// exists
func (d *Datastore) Query(id uint) (*kolide.Query, error) {
	sql := `
		SELECT * FROM queries WHERE id = ? AND NOT deleted
	`
	query := &kolide.Query{}
	if err := d.db.Get(query, sql, id); err != nil {
		return nil, err
	}

	return query, nil
}

// ListQueries returns a list of queries with sort order and results limit
// determined by passed in kolide.ListOptions
func (d *Datastore) ListQueries(opt kolide.ListOptions) ([]*kolide.Query, error) {
	sql := `
		SELECT * FROM queries WHERE NOT deleted
	`
	sql = appendListOptionsToSQL(sql, opt)
	results := []*kolide.Query{}

	if err := d.db.Select(&results, sql); err != nil {
		return nil, err
	}

	return results, nil

}
