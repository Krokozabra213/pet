package chatBusiness

import (
	"fmt"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/configs/chatconfig"

	// "github.com/Krokozabra213/sso/internal/chat/domain"
	// "github.com/Krokozabra213/sso/internal/chat/repository/broker"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

const (
	BufferSize = 100
)

type IClientRepo interface {
	Subscribe(client custombroker.IClient) error
	Unsubscribe(id int64) error
}

type IMessageRepo interface {
	SendMessage(msg string) error
}

type Chat struct {
	log        *slog.Logger
	cfg        *chatconfig.Config
	clientRepo IClientRepo
	msgRepo    IMessageRepo
}

func New(
	log *slog.Logger, cfg *chatconfig.Config, clientRepo IClientRepo,
	msgRepo IMessageRepo,
) *Chat {
	return &Chat{
		log:        log,
		cfg:        cfg,
		clientRepo: clientRepo,
		msgRepo:    msgRepo,
	}
}

func (a *Chat) Subscribe(username string) (<-chan interface{}, chan struct{}, uint64, error) {

	client := custombroker.NewClient(username, BufferSize)

	err := a.clientRepo.Subscribe(client)
	if err != nil {
		return nil, nil, 0, err
	}
	return client.GetBuffer(), client.GetDone(), client.GetUUID(), nil
}

func (a *Chat) Unsubscribe(uuid int64) error {

	a.clientRepo.Unsubscribe(userID)

	return nil
}

func (a *Chat) SendMessage(cliSendMsg *chat.ClientMessage_SendMessage) error {
	strMessage := cliSendMsg.SendMessage.GetContent()
	err := a.msgRepo.SendMessage(strMessage)
	if err != nil {
		return err
	}
	return nil
}

func (a *Chat) EntryPoint(clientMsg *chat.ClientMessage) error {

	switch msg := clientMsg.Type.(type) {
	case *chat.ClientMessage_SendMessage:
		err := a.SendMessage(msg)
		if err != nil {
			return err
		}
	case *chat.ClientMessage_Leave:
		err := a.Unsubscribe(msg.Leave.GetUserId())
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("wrong type message")
	}
	return nil
}
