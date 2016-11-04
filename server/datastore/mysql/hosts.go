package mysql

import (
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewHost(host *kolide.Host) (h *kolide.Host, e error) {
	return
}

func (d *Datastore) SaveHost(host *kolide.Host) (e error) {
	return
}

func (d *Datastore) DeleteHost(host *kolide.Host) (e error) {
	return
}

func (d *Datastore) Host(id uint) (h *kolide.Host, e error) {
	return
}

func (d *Datastore) ListHosts(opt kolide.ListOptions) (hosts []*kolide.Host, e error) {
	return
}

func (d *Datastore) EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (h *kolide.Host, e error) {
	return
}

func (d *Datastore) AuthenticateHost(nodeKey string) (h *kolide.Host, e error) {
	return
}

func (d *Datastore) MarkHostSeen(host *kolide.Host, t time.Time) (e error) {
	return
}

func (d *Datastore) SearchHosts(query string, omit []uint) (hosts []kolide.Host, e error) {
	return
}
