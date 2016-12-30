package kolide

import "fmt"

// Entity represents any object which has an unique identifier.
type Entity interface {
	EntityID() uint
}

// DBTable returns a database table which stores an entity.
func DBTable(e Entity) string {
	switch entity := e.(type) {
	case *Query:
		return "queries"
	case *Host:
		return "hosts"
	case *Label:
		return "labels"
	case *Pack:
		return "packs"
	case *Invite:
		return "invites"
	case *ScheduledQuery:
		return "scheduled_queries"
	default:
		// anyEntity is used to create a generic endpoint which can implement
		// the Deleter interface, as done in endpoint_delete.go
		type anyEntity interface {
			Entity
			EntityType() string
		}

		if e, ok := e.(anyEntity); ok {
			return e.EntityType()
		}

		panic(fmt.Sprintf(
			"entity %v has unknown db table: programmer error?", entity))
	}
}
