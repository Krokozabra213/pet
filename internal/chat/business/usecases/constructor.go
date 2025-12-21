package chatusecases

import (
	"log/slog"

	chatnewconfig "github.com/Krokozabra213/sso/newconfigs/chat"
)

const (
	BufferSize = 100
)

type Chat struct {
	log        *slog.Logger
	cfg        *chatnewconfig.Config
	clientRepo IClientRepo
	msgRepo    IMessageRepo
	msgSaver   IMessageSaver
}

func New(
	log *slog.Logger, cfg *chatnewconfig.Config, clientRepo IClientRepo,
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
