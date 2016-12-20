package inmem

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewOption(name string, optType kolide.OptionType, kolideRequires bool) (*kolide.Option, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	option := &kolide.Option{
		Name:              name,
		Type:              optType,
		RequiredForKolide: kolideRequires,
	}
	// don't allow dupe names
	for _, opt := range d.options {
		if opt.Name == name {
			return nil, fmt.Errorf("name '%s' is already in use", name)
		}
	}
	option.ID = d.nextID(option)
	d.options[option.ID] = option
	return option, nil
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

func (d *Datastore) SetOptionValues(vals []kolide.OptionValue) ([]kolide.OptionValue, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.optionValues = map[uint]*kolide.OptionValue{}
	for _, val := range vals {
		if val.ID == 0 {
			val.ID = d.nextID(val)
		}
		d.optionValues[val.ID] = &val
	}
	return vals, nil
}

func (d *Datastore) OptionValues() ([]kolide.OptionValue, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	response := []kolide.OptionValue{}
	for _, optVal := range d.optionValues {
		response = append(response, *optVal)
	}
	return response, nil
}
