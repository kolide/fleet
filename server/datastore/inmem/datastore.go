package inmem

import (
	"errors"
	"reflect"
	"sync"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/patrickmn/sortutil"
)

type Datastore struct {
	Driver  string
	mtx     sync.RWMutex
	nextIDs map[interface{}]uint

	users                map[uint]*kolide.User
	sessions             map[uint]*kolide.Session
	passwordResets       map[uint]*kolide.PasswordResetRequest
	invites              map[uint]*kolide.Invite
	labels               map[uint]*kolide.Label
	labelQueryExecutions map[uint]*kolide.LabelQueryExecution
	queries              map[uint]*kolide.Query
	packs                map[uint]*kolide.Pack
	hosts                map[uint]*kolide.Host
	packQueries          map[uint]*kolide.PackQuery
	packTargets          map[uint]*kolide.PackTarget

	orginfo *kolide.OrgInfo
}

func (orm *Datastore) Name() string {
	return "inmem"
}

func (orm *Datastore) Migrate() error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()
	orm.nextIDs = make(map[interface{}]uint)
	orm.users = make(map[uint]*kolide.User)
	orm.sessions = make(map[uint]*kolide.Session)
	orm.passwordResets = make(map[uint]*kolide.PasswordResetRequest)
	orm.invites = make(map[uint]*kolide.Invite)
	orm.labels = make(map[uint]*kolide.Label)
	orm.labelQueryExecutions = make(map[uint]*kolide.LabelQueryExecution)
	orm.queries = make(map[uint]*kolide.Query)
	orm.packs = make(map[uint]*kolide.Pack)
	orm.hosts = make(map[uint]*kolide.Host)
	orm.packQueries = make(map[uint]*kolide.PackQuery)
	orm.packTargets = make(map[uint]*kolide.PackTarget)
	return nil
}

func (orm *Datastore) Drop() error {
	return orm.Migrate()
}

// getLimitOffsetSliceBounds returns the bounds that should be used for
// re-slicing the results to comply with the requested ListOptions. Lack of
// generics forces us to do this rather than reslicing in this method.
func (orm *Datastore) getLimitOffsetSliceBounds(opt kolide.ListOptions, length int) (low uint, high uint) {
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

func sortResults(slice interface{}, opt kolide.ListOptions, fields map[string]string) error {
	field, ok := fields[opt.OrderKey]
	if !ok {
		return errors.New("cannot sort on unknown key: " + opt.OrderKey)
	}

	if opt.OrderDirection == kolide.OrderDescending {
		sortutil.DescByField(slice, field)
	} else {
		sortutil.AscByField(slice, field)
	}

	return nil
}

// nextID returns the next ID value that should be used for a struct of the
// given type
func (orm *Datastore) nextID(val interface{}) uint {
	valType := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(val)).Interface())
	orm.nextIDs[valType]++
	return orm.nextIDs[valType]
}
