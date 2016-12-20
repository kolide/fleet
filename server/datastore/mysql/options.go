package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewOption(name string, optType kolide.OptionType, kolideRequires bool) (*kolide.Option, error) {
	sqlStatement := `
    INSERT INTO options (
      name,
      type,
      required_for_kolide
    ) VALUES ( ?, ?, ? )
    `
	opt := &kolide.Option{
		Name:              name,
		Type:              optType,
		RequiredForKolide: kolideRequires,
	}
	result, err := d.db.Exec(sqlStatement, name, optType, kolideRequires)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	id, _ := result.LastInsertId()
	opt.ID = uint(id)
	return opt, nil
}

func (d *Datastore) Options() ([]kolide.Option, error) {
	sqlStatement := `
    SELECT *
    FROM options
    ORDER BY name ASC
  `
	var opts []kolide.Option
	if err := d.db.Select(&opts, sqlStatement); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return opts, nil
}

func (d *Datastore) SetOptionValues(opts []kolide.OptionValue) ([]kolide.OptionValue, error) {
	sqlStatement := `
    INSERT INTO option_values (
      option_id,
      value
    ) VALUES%s
    ON DUPLICATE KEY
    UPDATE
      value = VALUES(value),
			option_id = VALUES(option_id)
  `

	inList := []uint{}
	valuesClause := ""
	values := []interface{}{}
	for _, opt := range opts {
		if valuesClause != "" {
			valuesClause += ","
		}
		valuesClause += "(?, ?)"
		inList = append(inList, opt.OptionID)
		values = append(values, opt.OptionID, opt.Value)
	}
	sqlStatement = fmt.Sprintf(sqlStatement, valuesClause)
	_, err := d.db.Exec(sqlStatement, values...)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	// remove option values that weren't changed/created
	sqlStatement = `
    DELETE FROM option_values
    WHERE option_id NOT IN (?)
  `
	query, args, err := sqlx.In(sqlStatement, inList)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	query = d.db.Rebind(query)
	_, err = d.db.Exec(query, args...)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	results, err := d.OptionValues()
	if err != nil {
		return nil, err
	}

	return results, nil

}

func (d *Datastore) OptionValues() ([]kolide.OptionValue, error) {
	var results []kolide.OptionValue
	if err := d.db.Select(&results, "SELECT * FROM option_values"); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return results, nil
}
