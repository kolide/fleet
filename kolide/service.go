package kolide

type Service interface {
	UserService
}

func NewService(ds Datastore) Service {
	return service{
		// TODO set defaults
		bcryptCost:  10,
		saltKeySize: 10,
		db:          ds,
	}
}

type service struct {
	bcryptCost  int
	saltKeySize int
	db          Datastore
}
