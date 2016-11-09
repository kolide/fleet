package kolide

import "time"

// Createable contains common timestamp fields indicating create time
type CreateTimestamp struct {
	CreatedAt time.Time `json:"-" db:"created_at"`
}

// MarkAsCreated sets timestamp, intended to be called when Createable record is
// initially inserted in database
func (ci *CreateTimestamp) MarkAsCreated(created time.Time) {
	if ci.CreatedAt.IsZero() {
		ci.CreatedAt = created
	}
}

// Deleteable is used to indicate a record is deleted.  We don't actually
// delete record in the database. We mark it deleted, records with Deleted
// set to true will not normally be included in results
type DeleteFields struct {
	DeletedAt time.Time `json:"-" db:"deleted_at" gorm:"-"`
	Deleted   bool
}

// MarkDeleted indicates a record is deleted. It won't actually be removed from
// the database, but won't be returned in result sets.
func (d *DeleteFields) MarkDeleted(deleted time.Time) {
	d.DeletedAt = deleted
	d.Deleted = true
}

// UpdateTimestamp contains a timestamp that is set whenever an entity is changed
type UpdateTimestamp struct {
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// MarkAsUpdated is called when and entity is changed
func (u *UpdateTimestamp) MarkAsUpdated(updated time.Time) {
	u.UpdatedAt = updated
}

type UpdateCreateTimestamps struct {
	CreateTimestamp
	UpdateTimestamp
}

func (uct *UpdateCreateTimestamps) MarkAsCreated(created time.Time) {
	uct.CreateTimestamp.MarkAsCreated(created)
	uct.MarkAsUpdated(created)
}
