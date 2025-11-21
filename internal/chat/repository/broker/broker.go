package broker

import (
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
)

type Broker struct {
	B interface{}
}

func New() *Broker {
	return &Broker{}
}

func (br *Broker) Subscribe(client *domain.Client) error {
	// todo

	joinmessage := &chat.UserJoined{
		UserId:   client.ID,
		Username: client.Name,
	}
	client.Buf <- joinmessage

	sendmessage := &chat.ChatMessage{
		UserId:    client.ID,
		Username:  client.Name,
		Content:   "вы вошли в чат",
		Timestamp: time.Now().Unix(),
	}
	client.Buf <- sendmessage
	return nil
}

func (br *Broker) Unsubscribe(id int64) error {
	// todo
	return nil
}

func (br *Broker) SendMessage(msg string) error {
	// todo
	return nil
}

// type IClientRepo interface {
// 	Subscribe(client *domain.Client) error
// 	Unsubscribe(id int64) error
// }

// type IMessageRepo interface {
// 	SendMessage(msg string) error
// }
