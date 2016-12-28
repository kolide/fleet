package inmem

import (
	"reflect"
	"strings"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) Delete(e kolide.Entity) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	v := reflect.ValueOf(d).Elem().FieldByName(strings.Title(e.EntityType()))
	if v.IsValid() {
		v.SetMapIndex(reflect.ValueOf(e.EntityID()), reflect.Value{})
	}
	return nil
}
