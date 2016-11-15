package inmem

import (
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (orm *Inmem) NewAppConfig(info *kolide.AppConfig) (*kolide.AppConfig, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	info.ID = 1
	orm.orginfo = info
	return info, nil
}

func (orm *Inmem) AppConfig() (*kolide.AppConfig, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if orm.orginfo != nil {
		return orm.orginfo, nil
	}

	return nil, errors.ErrNotFound
}

func (orm *Inmem) SaveAppConfig(info *kolide.AppConfig) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	orm.orginfo = info
	return nil
}
