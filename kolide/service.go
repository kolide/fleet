package kolide

type Service interface {
	UserService
}

func NewService(ds Datastore) (Service, error) {
	return service{
		// TODO set defaults
		bcryptCost:  10,
		saltKeySize: 10,
		db:          ds,
	}, nil
}

type service struct {
	bcryptCost  int
	saltKeySize int
	db          Datastore
}
