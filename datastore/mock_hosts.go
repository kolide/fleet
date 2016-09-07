package datastore

import (
	"errors"

	"github.com/kolide/kolide-ose/kolide"
)

func (orm *mockDB) EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*kolide.Host, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()
	if uuid == "" {
		return nil, errors.New("missing uuid for host enrollment, programmer error?")
	}
	nodeKey, err := generateRandomText(nodeKeySize)
	if err != nil {
		return nil, err
	}
	host := &kolide.Host{
		UUID:      uuid,
		HostName:  hostname,
		IPAddress: ip,
		Platform:  platform,
		NodeKey:   nodeKey,
	}
	host.ID = uint(len(orm.hosts) + 1)
	orm.hosts[host.ID] = host

	return host, nil
}
