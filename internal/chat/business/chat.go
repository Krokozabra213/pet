package chatBusiness

import (
	"context"
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
	Subscribe(ctx context.Context, client custombroker.IClient) error
	Unsubscribe(uuid uint64) error
}

type IMessageRepo interface {
	SendMessage(ctx context.Context, message *chat.ClientMessage_SendMessage) error
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

func (a *Chat) Subscribe(ctx context.Context, username string) (<-chan interface{}, chan struct{}, uint64, error) {
	const op = "chat.Subscribe"
	log := a.log.With(
		slog.String("op", op),
	)

	client := custombroker.NewClient(username, BufferSize)
	err := a.clientRepo.Subscribe(ctx, client)
	if err != nil {
		log.Error("failed subscribe", "err", err)
		return nil, nil, 0, err
	}
	return client.GetBuffer(), client.GetDone(), client.GetUUID(), nil
}

func (a *Chat) Unsubscribe(uuid uint64) error {
	a.clientRepo.Unsubscribe(uuid)
	return nil
}

func (a *Chat) SendMessage(ctx context.Context, clientMsg *chat.ClientMessage) error {
	const op = "chat.SendMessage"
	log := a.log.With(
		slog.String("op", op),
	)

	switch msg := clientMsg.Type.(type) {
	case *chat.ClientMessage_SendMessage:
		err := a.msgRepo.SendMessage(ctx, msg)
		if err != nil {
			log.Error("failed send message", "err", err)
			return err
		}
	default:
		log.Error(fmt.Sprintf("wrong type message: %T", msg))
		return ErrWrongTypeMessage
	}
	return nil
}

// func (a *Chat) EntryPoint(clientMsg *chat.ClientMessage) error {

// 	switch msg := clientMsg.Type.(type) {
// 	case *chat.ClientMessage_SendMessage:
// 		err := a.SendMessage(msg)
// 		if err != nil {
// 			return err
// 		}
// 	case *chat.ClientMessage_Leave:
// 		err := a.Unsubscribe(msg.Leave.GetUserId())
// 		if err != nil {
// 			return err
// 		}
// 	default:
// 		return fmt.Errorf("wrong type message")
// 	}
// 	return nil
// }
