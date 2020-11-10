package mysql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

// NewUser creates a new user. If a user with the same username was
// soft-deleted, NewUser will replace the old one.
func (d *Datastore) NewUser(user *kolide.User) (*kolide.User, error) {
	var (
		deletedUser  kolide.User
		sqlStatement string
	)
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "begin NewUser transaction")
	}

	defer func() {
		if err != nil {
			rbErr := tx.Rollback()
			// It seems possible that there might be a case in
			// which the error we are dealing with here was thrown
			// by the call to tx.Commit(), and the docs suggest
			// this call would then result in sql.ErrTxDone.
			if rbErr != nil && rbErr != sql.ErrTxDone {
				panic(fmt.Sprintf("got err '%s' rolling back after err '%s'", rbErr, err))
			}
		}
	}()

	err = tx.Get(&deletedUser,
		"SELECT * FROM users WHERE username = ? AND deleted", user.Username)
	switch err {
	case nil:
		sqlStatement = `
			REPLACE INTO users (
				password,
				salt,
				name,
				username,
				email,
				admin,
				enabled,
				admin_forced_password_reset,
				gravatar_url,
				position,
				sso_enabled,
				deleted
			) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
		`
	case sql.ErrNoRows:
		sqlStatement = `
			INSERT INTO users (
				password,
				salt,
				name,
				username,
				email,
				admin,
				enabled,
				admin_forced_password_reset,
				gravatar_url,
				position,
				sso_enabled,
				deleted
			) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
		`
	default:
		return nil, errors.Wrap(err, "check for existing user")
	}
	deleted := false
	result, err := tx.Exec(sqlStatement, user.Password, user.Salt, user.Name,
		user.Username, user.Email, user.Admin, user.Enabled,
		user.AdminForcedPasswordReset, user.GravatarURL, user.Position,
		user.SSOEnabled, deleted)
	if err != nil && isDuplicate(err) {
		return nil, alreadyExists("User", deletedUser.ID)
	} else if err != nil {
		return nil, errors.Wrap(err, "create new user")
	}

	id, _ := result.LastInsertId()
	user.ID = uint(id)
	return user, nil
}

func (d *Datastore) findUser(searchCol string, searchVal interface{}) (*kolide.User, error) {
	sqlStatement := fmt.Sprintf(
		"SELECT * FROM users "+
			"WHERE %s = ? LIMIT 1",
		searchCol,
	)

	user := &kolide.User{}

	err := d.db.Get(user, sqlStatement, searchVal)
	if err != nil && err == sql.ErrNoRows {
		return nil, notFound("User").
			WithMessage(fmt.Sprintf("with %s=%v", searchCol, searchVal))
	} else if err != nil {
		return nil, errors.Wrap(err, "find user")
	}

	return user, nil
}

// User retrieves a user by name
func (d *Datastore) User(username string) (*kolide.User, error) {
	return d.findUser("username", username)
}

// ListUsers lists all users with limit, sort and offset passed in with
// kolide.ListOptions
func (d *Datastore) ListUsers(opt kolide.ListOptions) ([]*kolide.User, error) {
	sqlStatement := `
		SELECT * FROM users
	`
	sqlStatement = appendListOptionsToSQL(sqlStatement, opt)
	users := []*kolide.User{}

	if err := d.db.Select(&users, sqlStatement); err != nil {
		return nil, errors.Wrap(err, "list users")
	}

	return users, nil

}

func (d *Datastore) UserByEmail(email string) (*kolide.User, error) {
	return d.findUser("email", email)
}

func (d *Datastore) UserByID(id uint) (*kolide.User, error) {
	return d.findUser("id", id)
}

func (d *Datastore) SaveUser(user *kolide.User) error {
	sqlStatement := `
      UPDATE users SET
      	username = ?,
      	password = ?,
      	salt = ?,
      	name = ?,
      	email = ?,
      	admin = ?,
      	enabled = ?,
      	admin_forced_password_reset = ?,
      	gravatar_url = ?,
      	position = ?,
        sso_enabled = ?
      WHERE id = ?
      `
	result, err := d.db.Exec(sqlStatement, user.Username, user.Password,
		user.Salt, user.Name, user.Email, user.Admin, user.Enabled,
		user.AdminForcedPasswordReset, user.GravatarURL, user.Position, user.SSOEnabled, user.ID)
	if err != nil {
		return errors.Wrap(err, "save user")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "rows affected save user")
	}
	if rows == 0 {
		return notFound("User").WithID(user.ID)
	}

	return nil
}

// DeleteUserByID (soft) deletes the existing user object with the provided ID.
func (d *Datastore) DeleteUserByID(id uint) error {
	return d.deleteEntity("users", id)
}

// DeleteUsers (soft) deletes the existing user objects with the provided IDs.
// The number of deleted queries is returned along with any error.
func (d *Datastore) DeleteUsers(ids []uint) (uint, error) {
	sql := `
		UPDATE users
			SET deleted_at = NOW(), deleted = true
			WHERE id IN (?) AND NOT deleted
	`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return 0, errors.Wrap(err, "building delete users query")
	}

	result, err := d.db.Exec(query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "updating delete users query")
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "fetching delete users rows effected")
	}

	return uint(deleted), nil
}
