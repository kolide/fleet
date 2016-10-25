package kolide

import "context"

type QueryResultStore interface {
	// WriteResult writes a distributed query result submitted by an osqueryd client
	WriteResult(result DistributedQueryResult) error

	// ReadChannel returns a channel to be read for incoming distributed
	// query results
	ReadChannel(ctx context.Context, query DistributedQueryCampaign) (<-chan DistributedQueryResult, error)
}
