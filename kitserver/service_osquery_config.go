package kitserver

import (
	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/net/context"
)

func (svc service) GetAllQueries(ctx context.Context) ([]*kolide.Query, error) {
	return svc.ds.Queries()
}

func (svc service) GetQuery(ctx context.Context, id uint) (*kolide.Query, error) {
	return svc.ds.Query(id)
}

func (svc service) CreateQuery(ctx context.Context, p kolide.QueryPayload) error {
	return svc.ds.NewQuery(&kolide.Query{
		Name:         *p.Name,
		Query:        *p.Query,
		Interval:     *p.Interval,
		Snapshot:     *p.Snapshot,
		Differential: *p.Differential,
		Platform:     *p.Platform,
		Version:      *p.Version,
	})
}

func (svc service) ModifyQuery(ctx context.Context, id uint, p kolide.QueryPayload) (*kolide.Query, error) {
	query, err := svc.ds.Query(id)
	if err != nil {
		return nil, err
	}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	if p.Interval != nil {
		query.Interval = *p.Interval
	}

	if p.Snapshot != nil {
		query.Snapshot = *p.Snapshot
	}

	if p.Differential != nil {
		query.Differential = *p.Differential
	}

	if p.Platform != nil {
		query.Platform = *p.Platform
	}

	if p.Version != nil {
		query.Version = *p.Version
	}

	err = svc.ds.SaveQuery(query)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (svc service) DeleteQuery(ctx context.Context, id uint) error {
	query, err := svc.ds.Query(id)
	if err != nil {
		return err
	}

	err = svc.ds.DeleteQuery(query)
	if err != nil {
		return err
	}

	return nil
}

func (svc service) GetAllPacks(ctx context.Context) ([]*kolide.Pack, error) {
	return svc.ds.Packs()
}

func (svc service) GetPack(ctx context.Context, id uint) (*kolide.Pack, error) {
	return svc.ds.Pack(id)
}

func (svc service) CreatePack(ctx context.Context, p kolide.PackPayload) error {
	return svc.ds.NewPack(&kolide.Pack{
		Name:     *p.Name,
		Platform: *p.Platform,
	})
}

func (svc service) ModifyPack(ctx context.Context, id uint, p kolide.PackPayload) (*kolide.Pack, error) {
	pack, err := svc.ds.Pack(id)
	if err != nil {
		return nil, err
	}

	if p.Name != nil {
		pack.Name = *p.Name
	}

	if p.Platform != nil {
		pack.Platform = *p.Platform
	}

	err = svc.ds.SavePack(pack)
	if err != nil {
		return nil, err
	}

	return pack, err
}

func (svc service) DeletePack(ctx context.Context, id uint) error {
	pack, err := svc.ds.Pack(id)
	if err != nil {
		return err
	}

	err = svc.ds.DeletePack(pack)
	if err != nil {
		return err
	}

	return nil
}

func (svc service) AddQueryToPack(ctx context.Context, qid, pid uint) error {
	pack, err := svc.ds.Pack(pid)
	if err != nil {
		return err
	}

	query, err := svc.ds.Query(qid)
	if err != nil {
		return err
	}

	err = svc.ds.AddQueryToPack(query, pack)
	if err != nil {
		return err
	}

	return nil
}

func (svc service) RemoveQueryFromPack(ctx context.Context, qid, pid uint) error {
	pack, err := svc.ds.Pack(pid)
	if err != nil {
		return err
	}

	query, err := svc.ds.Query(qid)
	if err != nil {
		return err
	}

	err = svc.ds.RemoveQueryFromPack(query, pack)
	if err != nil {
		return err
	}

	return nil
}
