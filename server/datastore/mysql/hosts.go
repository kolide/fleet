package mysql

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewHost(host *kolide.Host) (*kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) SaveHost(host *kolide.Host) error {
	panic("not implemented")
}

func (d *Datastore) DeleteHost(host *kolide.Host) error {
	panic("not implemented")
}

func (d *Datastore) Host(id uint) (*kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) ListHosts(opt kolide.ListOptions) ([]*kolide.Host, error) {
	panic("not implemented")
}

// EnrollHost enrolls a host
func (d *Datastore) EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*kolide.Host, error) {
	if uuid == "" {
		return nil, errors.New("missing uuid for host enrollment", "programmer error")
	}
	// REVIEW If a deleted host is enrolled, it is undeleted
	sqlInsert := `
		INSERT INTO hosts (
			created_at,
			updated_at,
			detail_update_time,
			node_key,
			host_name,
			uuid,
			platform,
			primary_ip
		) VALUES (?, ?, ?, ?, ?, ?, ?, ? )
		ON DUPLICATE KEY UPDATE
			updated_at = VALUES(updated_at),
			detail_update_time = VALUES(detail_update_time),
			node_key = VALUES(node_key),
			host_name = VALUES(host_name),
			platform = VALUES(platform),
			primary_ip = VALUES(primary_ip),
			deleted = FALSE
	`
	args := []interface{}{}
	args = append(args, d.clock.Now())
	args = append(args, d.clock.Now())
	args = append(args, time.Unix(0, 0).Add(24*time.Hour))

	nodeKey, err := kolide.GenerateRandomText(nodeKeySize)

	args = append(args, nodeKey)
	args = append(args, hostname)
	args = append(args, uuid)
	args = append(args, platform)
	args = append(args, ip)

	var result sql.Result

	result, err = d.db.Exec(sqlInsert, args...)

	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	id, _ := result.LastInsertId()
	sqlSelect := `
		SELECT * FROM hosts WHERE id = ? LIMIT 1
	`
	host := &kolide.Host{}
	err = d.db.Get(host, sqlSelect, id)

	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	return host, nil

}

func (d *Datastore) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	sqlStatement := `
		SELECT * FROM hosts
		WHERE node_key = ? AND NOT DELETED LIMIT 1
	`
	host := &kolide.Host{}
	if err := d.db.Get(host, sqlStatement, nodeKey); err != nil {
		switch err {
		case sql.ErrNoRows:
			e := errors.NewFromError(err, http.StatusUnauthorized, "invalid node key")
			e.Extra = map[string]interface{}{"node_invalid": "true"}
			return nil, e
		default:
			return nil, errors.DatabaseError(err)
		}
	}

	return host, nil

}

func (d *Datastore) MarkHostSeen(*kolide.Host, time.Time) error {
	panic("not implemented")
}

func (d *Datastore) SearchHosts(query string, omit []uint) ([]kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) DistributedQueriesForHost(host *kolide.Host) (map[uint]string, error) {
	panic("not implmented")
}
