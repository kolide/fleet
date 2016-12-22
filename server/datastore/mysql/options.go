package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) SaveOption(opt kolide.Option) (*kolide.Option, error) {
	sqlStatement := `
    INSERT INTO options (
      name,
      type,
			value,
      read_only
    ) VALUES ( ?, ?, ?, ? )
		ON DUPLICATE KEY UPDATE
			value = VALUES(value)
    `
	result, err := d.db.Exec(
		sqlStatement,
		opt.Name,
		opt.Type,
		opt.RawValue,
		opt.ReadOnly,
	)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	if id != 0 {
		// assign id if we inserted a record
		opt.ID = uint(id)
	}
	return &opt, nil
}

func (d *Datastore) Option(id uint) (*kolide.Option, error) {
	sqlStatement := `
		SELECT *
		FROM options
		WHERE id = ?
	`
	var opt kolide.Option
	if err := d.db.Get(opt, sqlStatement, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Option").WithID(id)
		}
		return nil, errors.Wrap(err, "select option by ID")
	}
	return &opt, nil
}

func (d *Datastore) Options() ([]kolide.Option, error) {
	sqlStatement := `
    SELECT *
    FROM options
    ORDER BY name ASC
  `
	var opts []kolide.Option
	if err := d.db.Select(&opts, sqlStatement); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Option")
		}
		return nil, errors.Wrap(err, "select from options")
	}
	return opts, nil
}
