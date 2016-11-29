package inmem

import "fmt"

type notFoundError struct {
	ID           uint
	ResourceType string
}

func notFound(kind string, id uint) error {
	return &notFoundError{
		ID:           id,
		ResourceType: kind,
	}
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("%s %d was not found in the datastore", e.ResourceType, e.ID)
}

func (e *notFoundError) IsNotFound() bool {
	return true
}

type existsError struct {
	ID           uint
	ResourceType string
}

func alreadyExists(kind string, id uint) error {
	return &existsError{
		ID:           id,
		ResourceType: kind,
	}
}

func (e *existsError) Error() string {
	return fmt.Sprintf("%s %d already exists in the datastore", e.ResourceType, e.ID)
}

func (e *existsError) IsExists() bool {
	return true
}
