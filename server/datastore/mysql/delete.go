package mysql

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

// Delete updates a DB entity to mark it as deleted.
func (d *Datastore) Delete(e kolide.Entity) error {
	deleteStmt := fmt.Sprintf(
		`
	UPDATE %s SET deleted_at = NOW(), deleted = TRUE
		WHERE id = %d
		`,
		e.EntityType(), e.EntityID(),
	)
	result, err := d.db.Exec(deleteStmt)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete %s", e.EntityType()))
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		return notFound(e.EntityType()).WithID(e.EntityID())
	}
	return nil
}
