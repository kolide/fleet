package inmem

import "github.com/kolide/kolide-ose/server/kolide"

func (ds *Datastore) SaveLicense(string) error {
	panic("inmem is being deprecated")
}

func (ds *Datastore) License() (*kolide.License, error) {
	panic("inmem is being deprecated")
}
