package inmem

import "github.com/kolide/kolide-ose/server/kolide"

func (ds *Datastore) SaveLicense(token, key string) (*kolide.License, error) {
	panic("inmem is being deprecated")
}

func (ds *Datastore) License() (*kolide.License, error) {
	panic("inmem is being deprecated")
}

func (ds *Datastore) PublicKey(string) (string, error) {
	panic("inmem is being deprecated")
}
