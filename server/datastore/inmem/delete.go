package inmem

import (
	"reflect"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) Delete(e kolide.Entity) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	field := d.mapStructField(e)
	if field.IsValid() {
		field.SetMapIndex(reflect.ValueOf(e.EntityID()), reflect.Value{})
	}
	return nil
}
