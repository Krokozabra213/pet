package broker

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
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

	// for tests
	// joinmessage := &chat.UserJoined{
	// 	UserId:   int64(client.GetUUID()),
	// 	Username: client.GetName(),
	// }
	// client.Buf <- joinmessage

	// sendmessage := &chat.ChatMessage{
	// 	UserId:    client.ID,
	// 	Username:  client.Name,
	// 	Content:   "вы вошли в чат",
	// 	Timestamp: time.Now().Unix(),
	// }
	// client.Buf <- sendmessage
	return nil
}

func (br *Broker) Unsubscribe(uuid uint64) error {
	br.B.Unsubscribe(uuid)
	return nil
}

func (br *Broker) SendMessage(ctx context.Context, msg *chat.ClientMessage_SendMessage) error {
	err := br.B.Send(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
