package chatBusiness

import (
	"context"
	"log/slog"

	"github.com/Krokozabra213/sso/configs/chatconfig"
	"github.com/Krokozabra213/sso/internal/chat/domain"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"

	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

const (
	BufferSize = 100
)

type IClientRepo interface {
	Subscribe(ctx context.Context, client custombroker.IClient) error
	Unsubscribe(uuid uint64)
}

type IMessageRepo interface {
	Message(ctx context.Context, message interface{}) error
}

type IDefaultMessageSaver interface {
	SaveDefaultMessage(ctx context.Context, message *domain.DefaultMessage) (*domain.DefaultMessage, error)
}

type Chat struct {
	log             *slog.Logger
	cfg             *chatconfig.Config
	clientRepo      IClientRepo
	msgRepo         IMessageRepo
	defaultMsgSaver IDefaultMessageSaver
}

func New(
	log *slog.Logger, cfg *chatconfig.Config, clientRepo IClientRepo,
	msgRepo IMessageRepo, defaultMsgSaver IDefaultMessageSaver,
) *Chat {
	return &Chat{
		log:             log,
		cfg:             cfg,
		clientRepo:      clientRepo,
		msgRepo:         msgRepo,
		defaultMsgSaver: defaultMsgSaver,
	}
}

func (a *Chat) Subscribe(ctx context.Context, username string) (chatgrpc.IChatClient, error) {

	const op = "chat.Subscribe-Business"
	log := a.log.With(
		slog.String("op", op),
	)

	client := custombroker.NewClient(username, BufferSize)
	err := a.clientRepo.Subscribe(ctx, client)
	if err != nil {
		log.Error("failed subscribe", "err", err)
		return nil, err
	}
	return client, nil
}

func (a *Chat) Unsubscribe(uuid uint64) {
	a.clientRepo.Unsubscribe(uuid)
}

func (a *Chat) SendMessage(ctx context.Context, msg *domain.DefaultMessage) error {
	const op = "chat.SendMessage"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("message sended", "msg", msg)

	savedMsg, err := a.defaultMsgSaver.SaveDefaultMessage(ctx, msg)
	if err != nil {
		return err
	}

	err = a.msgRepo.Message(ctx, savedMsg)
	if err != nil {
		return err
	}

	return nil
}
