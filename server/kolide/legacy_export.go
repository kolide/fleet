package kolide

import "context"

// LegacyExportService interface describes methods that export the legacy
// (pre-fleetctl) configurations into fleetctl compatible implementations.
type LegacyExportService interface {
	// ExportConfig exports the Options and FIM configurations.
	ExportConfig(ctx context.Context) (string, error)
}
