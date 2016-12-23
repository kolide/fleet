package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) OptionByName(name string) (*kolide.Option, error) {
	sqlStatement := `
			SELECT *
			FROM options
			WHERE name = ?
		`
	var option kolide.Option
	if err := d.db.Get(&option, sqlStatement, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("option")
		}
		return nil, errors.Wrap(err, sqlStatement)
	}
	return &option, nil
}

func (d *Datastore) SaveOption(opt kolide.Option) error {
	var existing kolide.Option
	err := d.db.Get(&existing, "SELECT * FROM options WHERE id = ?", opt.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return notFound("option").WithID(opt.ID)
		}
		return errors.Wrap(err, "select from options")
	}
	// since we validate with passed in type verify that the passed
	// in type matches the type we have
	if existing.Type != opt.Type {
		return errors.New("type mismatch")
	}
	if existing.ReadOnly {
		return errors.New("readonly option can't be changed")
	}
	sqlStatement := `
    UPDATE options
		SET value = ?
		WHERE id = ?
	`
	_, err = d.db.Exec(
		sqlStatement,
		opt.RawValue,
		opt.ID,
	)
	if err != nil {
		return errors.Wrap(err, "update options")
	}
	return nil
}

func (d *Datastore) Option(id uint) (*kolide.Option, error) {
	sqlStatement := `
		SELECT *
		FROM options
		WHERE id = ?
	`
	var opt kolide.Option
	if err := d.db.Get(&opt, sqlStatement, id); err != nil {
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
