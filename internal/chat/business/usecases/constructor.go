package chatusecases

import (
	"log/slog"

	"github.com/Krokozabra213/sso/configs/chatconfig"
)

const (
	BufferSize = 100
)

type Chat struct {
	log        *slog.Logger
	cfg        *chatconfig.Config
	clientRepo IClientRepo
	msgRepo    IMessageRepo
	msgSaver   IMessageSaver
}

func New(
	log *slog.Logger, cfg *chatconfig.Config, clientRepo IClientRepo,
	msgRepo IMessageRepo, msgSaver IMessageSaver,
) *Chat {
	return &Chat{
		log:        log,
		cfg:        cfg,
		clientRepo: clientRepo,
		msgRepo:    msgRepo,
		msgSaver:   msgSaver,
	}
}
