package mysql

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kolide/kolide/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) CountHostsInTargets(hostIDs []uint, labelIDs []uint, now time.Time, onlineInterval time.Duration) (kolide.TargetMetrics, error) {
	sql := `
SELECT
COUNT(*) total,
COALESCE(SUM(CASE WHEN DATE_ADD(seen_time, INTERVAL 30 DAY) <= ? THEN 1 ELSE 0 END), 0) mia,
COALESCE(SUM(CASE WHEN DATE_ADD(seen_time, INTERVAL ? SECOND) <= ? AND DATE_ADD(seen_time, INTERVAL 30 DAY) >= ? THEN 1 ELSE 0 END), 0) offline,
COALESCE(SUM(CASE WHEN DATE_ADD(seen_time, INTERVAL ? SECOND) > ? THEN 1 ELSE 0 END), 0) online,
COALESCE(SUM(CASE WHEN DATE_ADD(created_at, INTERVAL 1 DAY) >= ? THEN 1 ELSE 0 END), 0) new
FROM
hosts h
		WHERE id IN (?)
OR (id IN (SELECT DISTINCT host_id FROM label_query_executions WHERE label_id IN (?) AND matches = 1))
		AND NOT deleted
`
	// DIRTY HACK -- FIX
	labelIDs = append(labelIDs, 0)
	hostIDs = append(hostIDs, 0)

	query, args, err := sqlx.In(sql, now, onlineInterval.Seconds(), now, now, onlineInterval.Seconds(), now, now, hostIDs, labelIDs)
	if err != nil {
		return kolide.TargetMetrics{}, errors.Wrap(err, "sqlx.In CountHostsInTargets")
	}

	res := []kolide.TargetMetrics{}
	err = d.db.Select(&res, query, args...)
	if err != nil {
		return kolide.TargetMetrics{}, errors.Wrap(err, "sqlx.Get CountHostsInTargets")
	}

	return res[0], nil
}
