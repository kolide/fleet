package kolide

type DecoratorStore interface {
	NewDecorator(decorator *Decorator) (*Decorator, error)
	DeleteDecorator(id uint) error
	Decorator(id uint) (*Decorator, error)
	ListDecorators() ([]*Decorator, error)
}

type DecoratorType int

const (
	DecoratorLoad DecoratorType = iota
	DecoratorAlways
	DecoratorInterval
)

type Decorator struct {
	UpdateCreateTimestamps
	ID       uint
	Type     DecoratorType
	Interval uint
	Query    string
}
