package inmem

import (
	"errors"
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) OptionByName(name string) (*kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	for _, opt := range d.options {
		if opt.Name == name {
			return opt, nil
		}
	}
	return nil, notFound("options")
}

func (d *Datastore) SaveOption(opt kolide.Option) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	existing, ok := d.options[opt.ID]
	if !ok {
		return notFound("option").WithID(opt.ID)
	}
	// since we will validate against the passed in type in the validation layer
	// we need to make sure that the passed in type matches the type we have
	if existing.Type != opt.Type {
		return fmt.Errorf("type mismatch")
	}
	if existing.ReadOnly {
		return errors.New("readonly option can't be changed")
	}
	if opt.RawValue == nil {
		existing.RawValue = nil
		return nil
	}

	existing.RawValue = new(string)
	*existing.RawValue = *opt.RawValue
	return nil
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
