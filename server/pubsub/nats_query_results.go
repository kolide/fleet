package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kolide/fleet/server/kolide"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type natsQueryResults struct {
	topicBase string
	conn      *nats.Conn
}

var _ kolide.QueryResultStore = &natsQueryResults{}

// NewNatsQueryResults creats a new Nats implementation of the
// QueryResultStore interface using the provided Nats connection pool.
func NewNatsQueryResults(conn *nats.Conn, topicBase string) *natsQueryResults {
	return &natsQueryResults{
		conn:      conn,
		topicBase: topicBase,
	}
}

func (nqr *natsQueryResults) pubSubForID(id uint) string {
	return fmt.Sprintf("%s.%d", nqr.topicBase, id)
}

func (nqr *natsQueryResults) WriteResult(result kolide.DistributedQueryResult) error {

	topic := nqr.pubSubForID(result.DistributedQueryCampaignID)

	jsonVal, err := json.Marshal(&result)
	if err != nil {
		return errors.Wrap(err, "marshalling JSON for result")
	}

	err = nqr.conn.Publish(topic, jsonVal)
	if err != nil {
		return errors.Wrap(err, "PUBLISH failed to Topic "+topic)
	}
	return nil
}

func (nqr *natsQueryResults) ReadChannel(ctx context.Context, query kolide.DistributedQueryCampaign) (<-chan interface{}, error) {
	outChannel := make(chan interface{})

	topic := nqr.pubSubForID(query.ID)

	// Channel Subscriber
	msgChannel := make(chan *nats.Msg, 64)
	sub, err := nqr.conn.ChanSubscribe(topic, msgChannel)
	if err != nil {
		return outChannel, err
	}

	go func() {
		defer sub.Unsubscribe()

		for {
			// Loop reading messages from conn.Receive() (via
			// msgChannel) until the context is cancelled.
			select {
			case msg, ok := <-msgChannel:
				if !ok {
					return
				}
				var res kolide.DistributedQueryResult
				err := json.Unmarshal(msg.Data, &res)
				if err != nil {
					outChannel <- err
				}
				outChannel <- res
			case <-ctx.Done():
				sub.Unsubscribe()
			}
		}

	}()
	return outChannel, nil
}

