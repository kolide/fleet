package inmem

import (
	"reflect"
	"strings"

	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (d *Datastore) Delete(ctx context.Context, e kolide.Entity) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	mapWords := func(all ...string) string {
		var out string
		for _, s := range all {
			out = out + strings.Title(s)
		}
		return out
	}
	// use reflect to get the value of a field which matches the EntityType()
	// for example , `invites` would get the datastore.Invites field,
	// which is a map. Setting the map key to reflect.Value{} achieves the same
	// result as delete(map, id) would.
	dbTable := mapWords(strings.Split(kolide.DBTable(e), "_")...)
	field := reflect.ValueOf(d).Elem().FieldByName(dbTable)
	if field.IsValid() {
		field.SetMapIndex(reflect.ValueOf(e.EntityID()), reflect.Value{})
	}
	return nil
}
