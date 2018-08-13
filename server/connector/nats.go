package connector

import (
	"github.com/nats-io/go-nats"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/health"
	"github.com/pkg/errors"
)

// NewNatsConnection create a new NatsConnection 
func NewNatsConn(conf config.NatsConfig) (*nats.Conn, error) {
	return nats.Connect(conf.URL)
}

type natsHealthChecker struct {
	conn *nats.Conn
}
var _ health.Checker = &natsHealthChecker{}

func NewNatsHealthChecker(conn *nats.Conn) (*natsHealthChecker, error) {
	return &natsHealthChecker{conn: conn}, nil
}

// HealthCheck verifies that the NATS  backend connected or connecting,  returning an
// error otherwise.
func (nhc *natsHealthChecker) HealthCheck() error {
	s := nhc.conn.Status()
	if s == nats.CONNECTED || s == nats.CONNECTING {
		return nil
	}
	return errors.New("Nats is not connected or connecting")
}

