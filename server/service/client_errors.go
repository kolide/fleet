package service

type SetupAlreadyErr interface {
	SetupAlready() bool
	Error() string
}

type setupAlreadyErr struct {
	reason string
}

func (e setupAlreadyErr) Error() string {
	return e.reason
}

func (e setupAlreadyErr) SetupAlready() bool {
	return true
}

func setupAlready() error {
	return setupAlreadyErr{
		reason: "Kolide Fleet has already been setup",
	}
}
