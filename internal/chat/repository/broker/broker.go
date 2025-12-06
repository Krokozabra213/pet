package broker

import (
	"context"

	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

type Broker struct {
	B *custombroker.CBroker
}

func New(brokerConn *custombroker.CBroker) *Broker {
	return &Broker{
		B: brokerConn,
	}
}

func (br *Broker) Subscribe(ctx context.Context, client custombroker.IClient) error {
	err := br.B.Subscribe(ctx, client)
	if err != nil {
		return err
	}

	return nil
}

func (br *Broker) Unsubscribe(uuid uint64) {
	br.B.Unsubscribe(uuid)
}

func (br *Broker) Message(ctx context.Context, msg interface{}) error {
	err := br.B.Send(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
