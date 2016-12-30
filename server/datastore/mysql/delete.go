package mysql

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Delete updates a DB entity to mark it as deleted.
func (d *Datastore) Delete(ctx context.Context, e kolide.Entity) error {
	dbTable := kolide.DBTable(e)
	deleteStmt := fmt.Sprintf(
		`
	UPDATE %s SET deleted_at = NOW(), deleted = TRUE
		WHERE id = %d
		`,
		dbTable, e.EntityID(),
	)
	result, err := d.db.Exec(deleteStmt)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete %s", dbTable))
	}
	rows, _ := result.RowsAffected()
	if rows != 1 {
		return notFound(dbTable).WithID(e.EntityID())
	}
	return nil
}
