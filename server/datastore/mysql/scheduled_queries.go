package mysql

import (
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) ListScheduledQueriesInPack(id uint, opts kolide.ListOptions) ([]*kolide.ScheduledQuery, error) {
	query := `
		SELECT sq.*, q.query, q.name
		FROM scheduled_queries sq
		JOIN queries q
		ON sq.query_id = q.id
		WHERE sq.pack_id = ?
		AND NOT sq.deleted
	`
	query = appendListOptionsToSQL(query, opts)
	results := []*kolide.ScheduledQuery{}

	if err := d.db.Select(&results, query, id); err != nil {
		return nil, errors.Wrap(err, "listing scheduled queries")
	}

	return results, nil
}
