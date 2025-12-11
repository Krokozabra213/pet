package chatusecases

import (
	"context"
	"log/slog"

	"github.com/Krokozabra213/sso/configs/chatconfig"
	"github.com/Krokozabra213/sso/internal/chat/domain"

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
