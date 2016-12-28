package service

type entity struct {
	Name string
	ID   uint
}

func (e *entity) EntityID() uint {
	return e.ID
}

func (e *entity) EntityType() string {
	return e.Name
}
