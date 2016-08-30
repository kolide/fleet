package kitserver

import "github.com/kolide/kolide-ose/kolide"

// configuration defaults
const (
	defaultBcryptCost  int    = 12
	defaultSaltKeySize int    = 24
	defaultCookieName  string = "KolideSession"
)

func NewService(ds kolide.Datastore) (kolide.Service, error) {
	var svc kolide.Service
	svc = service{
		bcryptCost:  defaultBcryptCost,
		saltKeySize: defaultSaltKeySize,
		ds:          ds,
	}
	svc = validationMiddleware{svc}
	return svc, nil
}

type service struct {
	bcryptCost  int
	saltKeySize int
	cookieName  string
	jwtKey      string
	ds          kolide.Datastore
}
