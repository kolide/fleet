package kolide

import "golang.org/x/net/context"

type QueryResultStore interface {
	// WriteResult writes a distributed query result submitted by an osqueryd client
	WriteResult(result DistributedQueryResult) error

	// ReadChannel returns a channel to be read for incoming distributed
	// query results. Channel values should be either
	// DistributedQueryResult or error
	ReadChannel(ctx context.Context, query DistributedQueryCampaign) (<-chan interface{}, error)
}
