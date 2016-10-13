package datastore

import (
	"sync"

	"github.com/kolide/kolide-ose/server/kolide"
)

type inmem struct {
	kolide.Datastore
	Driver string
	mtx    sync.RWMutex

	users                     map[uint]*kolide.User
	nextUserID                uint
	sessions                  map[uint]*kolide.Session
	nextSessionID             uint
	passwordResets            map[uint]*kolide.PasswordResetRequest
	nextPasswordResetID       uint
	invites                   map[uint]*kolide.Invite
	nextInviteID              uint
	labels                    map[uint]*kolide.Label
	nextLabelID               uint
	labelQueryExecutions      map[uint]*kolide.LabelQueryExecution
	nextLabelQueryExecutionID uint
	queries                   map[uint]*kolide.Query
	nextQueryID               uint
	packs                     map[uint]*kolide.Pack
	nextPackID                uint
	hosts                     map[uint]*kolide.Host
	nextHostID                uint

	orginfo *kolide.OrgInfo
}

func (orm *inmem) Name() string {
	return "inmem"
}

func (orm *inmem) Migrate() error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()
	orm.users = make(map[uint]*kolide.User)
	orm.sessions = make(map[uint]*kolide.Session)
	orm.passwordResets = make(map[uint]*kolide.PasswordResetRequest)
	orm.invites = make(map[uint]*kolide.Invite)
	orm.labels = make(map[uint]*kolide.Label)
	orm.labelQueryExecutions = make(map[uint]*kolide.LabelQueryExecution)
	orm.queries = make(map[uint]*kolide.Query)
	orm.packs = make(map[uint]*kolide.Pack)
	orm.hosts = make(map[uint]*kolide.Host)
	return nil
}

func (orm *inmem) Drop() error {
	return orm.Migrate()
}

// getLimitOffsetSliceBounds returns the bounds that should be used for
// re-slicing the results to comply with the requested ListOptions. Lack of
// generics forces us to do this rather than reslicing in this method.
func (orm *inmem) getLimitOffsetSliceBounds(opt kolide.ListOptions, length int) (low uint, high uint) {
	if opt.PerPage == 0 {
		// PerPage value of 0 indicates unlimited
		return 0, uint(length)
	}

	offset := opt.Page * opt.PerPage
	max := offset + opt.PerPage
	if offset > uint(length) {
		offset = uint(length)
	}
	if max > uint(length) {
		max = uint(length)
	}
	return offset, max
}
