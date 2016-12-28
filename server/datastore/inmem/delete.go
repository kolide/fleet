package inmem

import (
	"reflect"
	"strings"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) Delete(e kolide.Entity) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	// use reflect to get the value of a field which matches the EntityType()
	// for example , `invites` would get the datastore.Invites field,
	// which is a map. Setting the map key to reflect.Value{} achieves the same
	// result as delete(map, id) would.
	field := reflect.ValueOf(d).Elem().FieldByName(strings.Title(e.EntityType()))
	if field.IsValid() {
		field.SetMapIndex(reflect.ValueOf(e.EntityID()), reflect.Value{})
	}
	return nil
}
