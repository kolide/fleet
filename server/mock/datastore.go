package mock

import (
	"errors"
	"reflect"
	"time"

	"github.com/kolide/kolide/server/kolide"
)

type data struct {
	nextIDs map[interface{}]uint

	users              map[uint]*kolide.User
	sessions           map[uint]*kolide.Session
	options            map[uint]*kolide.Option
	packs              map[uint]*kolide.Pack
	packTargets        map[uint]*kolide.PackTarget
	labels             map[uint]*kolide.Label
	queries            map[uint]*kolide.Query
	scheduledQueries   map[uint]*kolide.ScheduledQuery
	decorators         map[uint]*kolide.Decorator
	fimSections        map[uint]*kolide.FIMSection
	yaraSignatureGroup map[uint]*kolide.YARASignatureGroup

	config *kolide.AppConfig
}

type operation func(d *data)

type Datastore struct {
	Store
	ops chan operation
}

func NewDatastore() *Datastore {

	ds := new(Datastore)

	ds.ops = make(chan operation)

	go func(ops <-chan operation) {
		d := new(data)
		d.nextIDs = make(map[interface{}]uint)
		d.users = make(map[uint]*kolide.User)
		d.sessions = make(map[uint]*kolide.Session)
		d.packs = make(map[uint]*kolide.Pack)
		d.labels = make(map[uint]*kolide.Label)
		d.packTargets = make(map[uint]*kolide.PackTarget)
		d.options = make(map[uint]*kolide.Option)
		d.queries = make(map[uint]*kolide.Query)
		d.scheduledQueries = make(map[uint]*kolide.ScheduledQuery)
		d.decorators = make(map[uint]*kolide.Decorator)
		d.fimSections = make(map[uint]*kolide.FIMSection)
		d.yaraSignatureGroup = make(map[uint]*kolide.YARASignatureGroup)

		for operateOn := range ops {
			operateOn(d)
		}
	}(ds.ops)

	ds.NewYARAFilePathFunc = func(fileSectionName, sigGroupName string) error {
		//res := make(chan struct{})
		return nil
	}

	ds.NewYARASignatureGroupFunc = func(psg *kolide.YARASignatureGroup) (*kolide.YARASignatureGroup, error) {
		res := make(chan uint)
		sg := *psg
		ds.ops <- func(d *data) {
			sg.ID = d.nextID(sg)
			d.yaraSignatureGroup[sg.ID] = &sg
			res <- sg.ID
		}
		psg.ID = <-res
		return psg, nil
	}

	ds.NewFIMSectionFunc = func(pfs *kolide.FIMSection) (*kolide.FIMSection, error) {
		res := make(chan uint)
		fs := *pfs
		ds.ops <- func(d *data) {
			fs.ID = d.nextID(fs)
			d.fimSections[fs.ID] = &fs
			res <- fs.ID
		}
		pfs.ID = <-res
		return pfs, nil
	}

	ds.QueryByNameFunc = func(name string) (*kolide.Query, bool, error) {
		res := make(chan interface{})
		ds.ops <- func(d *data) {
			for _, q := range d.queries {
				if q.Name == name {
					res <- q
					return
				}
			}
			res <- false
		}
		v := <-res
		if q, ok := v.(*kolide.Query); ok {
			return q, true, nil
		}
		return nil, false, nil
	}

	ds.NewQueryFunc = func(pq *kolide.Query) (*kolide.Query, error) {
		res := make(chan uint)
		q := *pq
		ds.ops <- func(d *data) {
			q.ID = d.nextID(q)
			d.queries[q.ID] = &q
			res <- q.ID
		}
		pq.ID = <-res
		return pq, nil
	}

	ds.NewScheduledQueryFunc = func(psq *kolide.ScheduledQuery) (*kolide.ScheduledQuery, error) {
		res := make(chan uint)
		sq := *psq
		ds.ops <- func(d *data) {
			sq.ID = d.nextID(sq)
			q := d.queries[sq.QueryID]
			sq.Query = q.Query
			sq.Name = q.Name
			d.scheduledQueries[sq.ID] = &sq
			res <- sq.ID
		}
		psq.ID = <-res
		return psq, nil
	}

	ds.NewDecoratorFunc = func(dec *kolide.Decorator) (*kolide.Decorator, error) {
		res := make(chan uint)
		dd := *dec
		ds.ops <- func(d *data) {
			dd.ID = d.nextID(dd)
			d.decorators[dd.ID] = &dd
			res <- dd.ID
		}
		dec.ID = <-res
		return dec, nil
	}

	ds.ListDecoratorsFunc = func() ([]*kolide.Decorator, error) {
		res := make(chan []*kolide.Decorator)
		ds.ops <- func(d *data) {
			var decs []*kolide.Decorator
			for _, dec := range d.decorators {
				decs = append(decs, dec)
			}
			res <- decs
		}
		return <-res, nil
	}

	ds.DeleteDecoratorFunc = func(id uint) error {
		res := make(chan interface{})
		ds.ops <- func(d *data) {
			if dec, ok := d.decorators[id]; ok {
				delete(d.decorators, id)
				res <- dec
				return
			}
			res <- errors.New("decorator not found")
		}
		v := <-res
		if _, ok := v.(*kolide.Decorator); ok {
			return nil
		}
		return v.(error)
	}

	ds.DecoratorFunc = func(id uint) (*kolide.Decorator, error) {
		res := make(chan interface{})
		ds.ops <- func(d *data) {
			if d, ok := d.decorators[id]; ok {
				res <- d
				return
			}
			res <- errors.New("decorator not found")
		}
		v := <-res
		if dec, ok := v.(*kolide.Decorator); ok {
			return dec, nil
		}
		return nil, v.(error)
	}

	ds.ListOptionsFunc = func() ([]kolide.Option, error) {
		res := make(chan []kolide.Option)
		ds.ops <- func(d *data) {
			var opts []kolide.Option
			for _, o := range d.options {
				opts = append(opts, *o)
			}
			res <- opts
		}
		return <-res, nil
	}

	ds.SaveOptionsFunc = func(opts []kolide.Option) error {
		res := make(chan struct{})
		ds.ops <- func(d *data) {
			for _, o := range opts {
				o.ID = d.nextID(&o)
				d.options[o.ID] = &o
			}
			res <- struct{}{}
		}
		<-res
		return nil
	}

	ds.OptionByNameFunc = func(name string) (*kolide.Option, error) {
		res := make(chan interface{})
		ds.ops <- func(d *data) {
			for _, o := range d.options {
				if o.Name == name {
					res <- o
					return
				}
			}
			res <- errors.New("no option with name " + name)
		}
		v := <-res
		if o, ok := v.(*kolide.Option); ok {
			return o, nil
		}
		return nil, v.(error)
	}

	ds.NewAppConfigFunc = func(c *kolide.AppConfig) (*kolide.AppConfig, error) {
		r := make(chan struct{})
		ds.ops <- func(d *data) {
			d.config = c
			r <- struct{}{}
		}
		<-r
		return c, nil
	}

	ds.AppConfigFunc = func() (*kolide.AppConfig, error) {
		r := make(chan kolide.AppConfig)
		ds.ops <- func(d *data) {
			r <- *d.config
		}
		c := <-r
		return &c, nil
	}

	ds.SaveAppConfigFunc = func(c *kolide.AppConfig) error {
		r := make(chan struct{})
		ds.ops <- func(d *data) {
			d.config = c
			r <- struct{}{}
		}
		<-r
		return nil
	}

	ds.NewUserFunc = func(pu *kolide.User) (*kolide.User, error) {
		r := make(chan uint)
		u := *pu
		ds.ops <- func(d *data) {
			u.ID = d.nextID(u)
			d.users[u.ID] = &u
			r <- u.ID
		}
		pu.ID = <-r
		return pu, nil
	}

	ds.UserFunc = func(username string) (*kolide.User, error) {
		r := make(chan interface{})
		ds.ops <- func(d *data) {
			for _, u := range d.users {
				if u.Username == username {
					r <- u
					return
				}
			}
			r <- errors.New("no user")
		}
		v := <-r
		if u, ok := v.(*kolide.User); ok {
			return u, nil
		}
		return nil, v.(error)
	}

	ds.UserByIDFunc = func(id uint) (*kolide.User, error) {
		r := make(chan interface{})
		ds.ops <- func(d *data) {
			if u, ok := d.users[id]; ok {
				r <- u
				return
			}
			r <- errors.New("no user")
		}
		v := <-r
		if u, ok := v.(*kolide.User); ok {
			return u, nil
		}
		return nil, v.(error)
	}

	ds.SaveUserFunc = func(u *kolide.User) error {
		res := make(chan interface{})
		ds.ops <- func(d *data) {
			if _, ok := d.users[u.ID]; ok {
				d.users[u.ID] = u
				res <- u
				return
			}
			res <- errors.New("no user")
		}
		v := <-res
		if _, ok := v.(*kolide.User); ok {
			return nil
		}
		return v.(error)
	}

	ds.NewSessionFunc = func(psess *kolide.Session) (*kolide.Session, error) {
		r := make(chan uint)
		sess := *psess
		ds.ops <- func(d *data) {
			sess.ID = d.nextID(sess)
			d.sessions[sess.ID] = &sess
			r <- sess.ID
		}
		psess.ID = <-r
		return psess, nil
	}

	ds.SessionByKeyFunc = func(k string) (*kolide.Session, error) {
		r := make(chan interface{})
		ds.ops <- func(d *data) {
			for _, sess := range d.sessions {
				if sess.Key == k {
					r <- sess
					return
				}
			}
			r <- errors.New("no such session")
		}
		v := <-r
		if s, ok := v.(*kolide.Session); ok {
			return s, nil
		}
		return nil, v.(error)
	}

	ds.MarkSessionAccessedFunc = func(sess *kolide.Session) error {
		r := make(chan interface{})
		ds.ops <- func(d *data) {
			if s, ok := d.sessions[sess.ID]; ok {
				s.AccessedAt = time.Now().UTC()
				r <- s
				return
			}
			r <- errors.New("no such session")
		}
		v := <-r
		if e, ok := v.(error); ok {
			return e
		}
		return nil
	}

	ds.PackByNameFunc = func(name string) (*kolide.Pack, bool, error) {
		res := make(chan interface{})
		ds.ops <- func(d *data) {
			for _, p := range d.packs {
				if p.Name == name {
					res <- p
					return
				}
			}
			res <- false
		}
		v := <-res
		if p, ok := v.(*kolide.Pack); ok {
			return p, true, nil
		}
		return nil, false, nil
	}

	ds.NewPackFunc = func(pp *kolide.Pack) (*kolide.Pack, error) {
		res := make(chan uint)
		p := *pp
		ds.ops <- func(d *data) {
			p.ID = d.nextID(p)
			d.packs[p.ID] = &p
			res <- p.ID
		}
		pp.ID = <-res
		return pp, nil
	}

	ds.NewLabelFunc = func(pl *kolide.Label) (*kolide.Label, error) {
		res := make(chan uint)
		l := *pl
		ds.ops <- func(d *data) {
			l.ID = d.nextID(l)
			d.labels[l.ID] = &l
			res <- l.ID
		}
		pl.ID = <-res
		return pl, nil
	}

	ds.AddLabelToPackFunc = func(lid, pid uint) error {
		res := make(chan struct{})
		ds.ops <- func(d *data) {
			for _, pt := range d.packTargets {
				if pt.PackID == pid && pt.Target.Type == kolide.TargetLabel && pt.TargetID == lid {
					res <- struct{}{}
					return
				}
			}
			pt := &kolide.PackTarget{
				PackID: pid,
				Target: kolide.Target{
					Type:     kolide.TargetLabel,
					TargetID: lid,
				},
			}
			pt.ID = d.nextID(pt)
			d.packTargets[pt.ID] = pt
			res <- struct{}{}
		}
		<-res
		return nil
	}

	return ds
}

func (ds *Datastore) Close() {
	close(ds.ops)
}

func (d *data) nextID(val interface{}) uint {
	valType := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(val)).Interface())
	d.nextIDs[valType]++
	return d.nextIDs[valType]
}
