package mysql

import (
	"time"

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

func (d *Datastore) EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) MarkHostSeen(*kolide.Host, time.Time) error {
	panic("not implemented")
}

func (d *Datastore) SearchHosts(query string, omit []uint) ([]kolide.Host, error) {
	panic("not implemented")
}
