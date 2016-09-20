package kolide

import (
	"time"

	"golang.org/x/net/context"
)

type HostStore interface {
	NewHost(host *Host) (*Host, error)
	SaveHost(host *Host) error
	DeleteHost(host *Host) error
	Host(id uint) (*Host, error)
	Hosts() ([]*Host, error)
	EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*Host, error)
	AuthenticateHost(nodeKey string) (*Host, error)
	MarkHostSeen(host *Host, t time.Time) error
}

type HostService interface {
	GetAllHosts(ctx context.Context) ([]*Host, error)
	GetHost(ctx context.Context, id uint) (*Host, error)
	NewHost(ctx context.Context, p HostPayload) (*Host, error)
	ModifyHost(ctx context.Context, id uint, p HostPayload) (*Host, error)
	DeleteHost(ctx context.Context, id uint) error
}

type HostPayload struct {
	NodeKey   *string
	HostName  *string
	UUID      *string
	IPAddress *string
	Platform  *string
}

type Host struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeKey   string `gorm:"unique_index:idx_host_unique_nodekey"`
	HostName  string
	UUID      string `gorm:"unique_index:idx_host_unique_uuid"`
	IPAddress string
	Platform  string
}

// NeedsDetailUpdate determines whether the host context (IP, platform, etc.)
// needs to be updated
func (h Host) NeedsDetailUpdate() bool {
	// Currently we only attempt to update platform
	return h.Platform == ""
}

const hostDetailQueryPrefix = "kolide_detail_query_"
const hostLabelQueryPrefix = "kolide_label_query_"

// GetDetailQueries returns the map of queries that should be executed by
// osqueryd to fill in the host details
func (h Host) GetDetailQueries() map[string]string {
	return map[string]string{
		hostDetailQueryPrefix + "platform": "select build_platform from osquery_info;",
	}
}
