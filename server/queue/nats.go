package queue

import (
	"encoding/json"
	"text/template"
	"bytes"
	"bufio"
	"fmt"
	"context"

	"github.com/nats-io/go-nats"
	kitlog "github.com/go-kit/kit/log"
	hostctx "github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/config"
)


type NatsQueue struct {
	logger        kitlog.Logger
	topicRaw      string
	topicTemplate *template.Template
	conn          *nats.Conn
}

func NewNatsQueue(appLogger kitlog.Logger, conn *nats.Conn, conf config.NatsQueueConfig) (*NatsQueue, error) {
	nq := &NatsQueue{
		conn:      conn,
		topicRaw:  conf.Topic,
	}
	t, err := template.New("nats_topic").Parse(nq.topicRaw)
	if err != nil {
		return nq, err
	}
	nq.topicTemplate = t
	return nq, nil
}

func (nq *NatsQueue) Messages(ctx context.Context, logs []json.RawMessage) error {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return fmt.Errorf("Context was not able to decode the host infomation")
	}

	var b bytes.Buffer // A Buffer needs no initialization.
	writer := bufio.NewWriter(&b)
	err := nq.topicTemplate.Execute(writer, host)
	if err != nil {
		return err 
	}
	for _, log := range logs {
		nq.conn.Publish(b.String(), log)
	}
	return nil 
}
