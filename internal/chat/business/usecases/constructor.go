package chatusecases

import (
	chatnewconfig "github.com/Krokozabra213/sso/newconfigs/chat"
)

const (
	BufferSize = 100
)

type Chat struct {
	cfg        *chatnewconfig.Config
	clientRepo IClientRepo
	msgRepo    IMessageRepo
	msgSaver   IMessageSaver
}

func New(
	cfg *chatnewconfig.Config, clientRepo IClientRepo,
	msgRepo IMessageRepo, msgSaver IMessageSaver,
) *Chat {
	return &Chat{
		cfg:        cfg,
		clientRepo: clientRepo,
		msgRepo:    msgRepo,
		msgSaver:   msgSaver,
	}
}
