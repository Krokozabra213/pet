package chatBusiness

import (
	"fmt"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/configs/chatconfig"
	"github.com/Krokozabra213/sso/internal/chat/domain"
)

type IClientRepo interface {
	Subscribe(id int64, name string, buffer chan interface{}) error
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

func (a *Chat) Subscribe(cli *domain.Client) error {
	// TODO
	return nil
}

func (a *Chat) EntryPoint(clientMsg *chat.ClientMessage) error {

	switch msg := clientMsg.Type.(type) {
	case *chat.ClientMessage_SendMessage:
		//todo

		err := a.SendMessage(msg)
		if err != nil {
			return err
		}
	case *chat.ClientMessage_Leave:
		//todo
		err := a.ClientLeave(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Chat) SendMessage(cliSendMsg *chat.ClientMessage_SendMessage) error {
	// TODO
	strMessage := cliSendMsg.SendMessage.GetContent()
	fmt.Println(strMessage)
	return nil
}

func (a *Chat) ClientLeave(cliLeaveMsg *chat.ClientMessage_Leave) error {
	// TODO
	userID := cliLeaveMsg.Leave.GetUserId()
	fmt.Println(userID)
	return nil
}

func (a *Chat) Unsubscribe(userID int64) error {
	// TODO
	return nil
}
