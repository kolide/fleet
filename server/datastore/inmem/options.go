package inmem

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) SaveOption(opt kolide.Option) (*kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if opt.ID == 0 {
		// don't allow dupe names
		for _, o := range d.options {
			if opt.Name == o.Name {
				return nil, fmt.Errorf("name '%s' is already in use", opt.Name)
			}
		}
		opt.ID = d.nextID(opt)
		d.options[opt.ID] = &opt
		return &opt, nil
	}

	saved := d.options[opt.ID]
	saved.RawValue = opt.RawValue
	return saved, nil

}

func (d *Datastore) Option(id uint) (*kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	saved, ok := d.options[id]
	if !ok {
		return nil, notFound("Option").WithID(id)
	}
	return saved, nil
}

func (d *Datastore) Options() ([]kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	result := []kolide.Option{}
	for _, opt := range d.options {
		result = append(result, *opt)
	}
	return result, nil
}
