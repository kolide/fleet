package mysql

import (
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewPasswordResetRequest(req *kolide.PasswordResetRequest) (*kolide.PasswordResetRequest, error) {
	req.MarkAsCreated(d.clock.Now())
	sqlStatement := `
		INSERT INTO password_reset_requests
		(created_at, updated_at, user_id, token)
		VALUES (?,?,?,?)
	`
	response, err := d.db.Exec(sqlStatement, req.CreatedAt, req.UpdatedAt, req.UserID,
		req.Token)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	id, _ := response.LastInsertId()
	req.ID = uint(id)
	return req, nil

}

func (d *Datastore) SavePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	panic("not implemented")
}

func (d *Datastore) DeletePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	panic("not implemented")
}

func (d *Datastore) DeletePasswordResetRequestsForUser(userID uint) error {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetByID(id uint) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetsByUserID(id uint) ([]*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetByToken(token string) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetByTokenAndUserID(token string, id uint) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}
