package kolide

type QueryResultStore interface {
	// WriteResult writes a distributed query result submitted by an osqueryd client
	WriteResult(result DistributedQueryResult) error

	// ReadChannel returns a channel to be read for incoming distributed
	// query results
	ReadChannel(query DistributedQueryCampaign) (chan DistributedQueryResult, error)

	// CloseQuery indicates that no more results will be read for the given
	// query campaign
	CloseQuery(query DistributedQueryCampaign)
}
