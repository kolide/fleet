package mysql

import (
	"database/sql"
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
	_, err := d.db.Exec(deleteStmt)
	if err != nil && err == sql.ErrNoRows {
		return notFound(e.EntityType()).WithID(e.EntityID())
	} else if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete %s", e.EntityType()))
	}
	return nil
}
