package datastore

import (
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (orm gormDB) NewQuery(query *kolide.Query) (*kolide.Query, error) {
	if query == nil {
		return nil, errors.New(
			"error creating query",
			"nil pointer passed to NewQuery",
		)
	}
	err := orm.DB.Create(query).Error
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (orm gormDB) SaveQuery(query *kolide.Query) error {
	if query == nil {
		return errors.New(
			"error saving query",
			"nil pointer passed to SaveQuery",
		)
	}
	return orm.DB.Save(query).Error
}

func (orm gormDB) DeleteQuery(query *kolide.Query) error {
	if query == nil {
		return errors.New(
			"error deleting query",
			"nil pointer passed to DeleteQuery",
		)
	}
	return orm.DB.Delete(query).Error
}

func (orm gormDB) Query(id uint) (*kolide.Query, error) {
	query := &kolide.Query{
		ID: id,
	}
	err := orm.DB.Where(query).First(query).Error
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (orm gormDB) ListQueries(opt kolide.ListOptions) ([]*kolide.Query, error) {
	var queries []*kolide.Query
	err := orm.applyListOptions(opt).Find(&queries).Error
	return queries, err
}

func (orm gormDB) DistributedQueriesForHost(host *kolide.Host) ([]kolide.Query, error) {
	sql := `
SELECT DISTINCT dqc.id, q.query
FROM distributed_query_campaigns dqc
JOIN distributed_query_campaign_targets dqct
    ON (dqc.id = dqct.distributed_query_campaign_id)
LEFT JOIN label_query_executions lqe
    ON (dqct.type = 0 AND dqct.target_id = lqe.label_id)
LEFT JOIN hosts h
    ON ((dqct.type = 0 AND lqe.host_id = h.id AND lqe.matches) OR (dqct.type = 1 AND dqct.target_id = h.id))
LEFT JOIN distributed_query_executions dqe
    ON (h.id = dqe.host_id AND dqc.id = dqe.distributed_query_id)
JOIN queries q
    ON (dqc.query_id = q.id)
WHERE dqe.status IS NULL AND h.id = 2;
`
	_ = sql
	return nil, nil
}
