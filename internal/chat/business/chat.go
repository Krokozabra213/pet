package chatBusiness

import (
	"fmt"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/configs/chatconfig"
	"github.com/Krokozabra213/sso/internal/chat/domain"
)

const (
	BufferSize = 100
)

type IClientRepo interface {
	Subscribe(client *domain.Client) error
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

func (a *Chat) Subscribe(userID int64, username string) (chan interface{}, chan struct{}, error) {

	client := &domain.Client{
		ID:   userID,
		Name: username,
		Buf:  make(chan interface{}, BufferSize), // для получения сообщений
		Done: make(chan struct{}),                // для graceful shutdown
	}

	fmt.Println(client)

	err := a.clientRepo.Subscribe(client)
	if err != nil {
		return nil, nil, err
	}

	return client.Buf, client.Done, nil
}

func (a *Chat) Unsubscribe(userID int64) error {

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
