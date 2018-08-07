package kolide


type HealthService interface {
	HealthCheckers() (map[string]interface{})
}
