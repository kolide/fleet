package inmem

import (
	"reflect"
	"strings"

	"github.com/kolide/kolide-ose/server/kolide"
)

// mapStructField returns the struct field that stores the entity indexed by ID
func (d *Datastore) mapStructField(e kolide.Entity) reflect.Value {
	mapWords := func(all ...string) string {
		var out string
		for _, s := range all {
			out = out + strings.Title(s)
		}
		return out
	}
	dbTable := mapWords(strings.Split(kolide.DBTable(e), "_")...)
	field := reflect.ValueOf(d).Elem().FieldByName(dbTable)
	return field
}

// byID is a helper that retrieves an entity from the Datastore by it's ID field
func (d *Datastore) byID(e kolide.Entity) (interface{}, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	field := d.mapStructField(e)
	v := field.MapIndex(reflect.ValueOf(e.EntityID()))
	if !v.IsValid() {
		return nil, notFound(kolide.DBTable(e)).WithID(e.EntityID())
	}
	return v.Interface(), nil
}
