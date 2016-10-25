package datastore

import (
	"sync"

	"github.com/kolide/kolide-ose/server/kolide"
)

type inmemQueryResults struct {
	resultChannels map[uint]chan kolide.DistributedQueryResult
	channelMutex   sync.Mutex
}

var _ kolide.QueryResultStore = &inmemQueryResults{}

func newInmemQueryResults() inmemQueryResults {
	return inmemQueryResults{resultChannels: map[uint]chan kolide.DistributedQueryResult{}}
}

func (im *inmemQueryResults) getChannel(id uint) chan kolide.DistributedQueryResult {
	im.channelMutex.Lock()
	defer im.channelMutex.Unlock()

	channel, ok := im.resultChannels[id]
	if !ok {
		channel = make(chan kolide.DistributedQueryResult)
		im.resultChannels[id] = channel
	}
	return channel
}

func (im *inmemQueryResults) WriteResult(result kolide.DistributedQueryResult) error {
	// Write the result
	im.getChannel(result.DistributedQueryCampaignID) <- result

	return nil
}

func (im *inmemQueryResults) ReadChannel(query kolide.DistributedQueryCampaign) (<-chan kolide.DistributedQueryResult, error) {
	return im.getChannel(query.ID), nil
}

func (im *inmemQueryResults) CloseQuery(query kolide.DistributedQueryCampaign) {
	channel, ok := im.resultChannels[query.ID]
	if !ok {
		return
	}
	close(channel)
}
