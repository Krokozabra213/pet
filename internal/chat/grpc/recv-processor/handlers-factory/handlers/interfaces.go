package handlers

import "github.com/Krokozabra213/protos/gen/go/chat"

const (
	ImageType uint8 = 1
	TextType  uint8 = 2
)

type IMessageHandler interface {
	GetType() string
	CanHandle(message *chat.ClientMessage) bool
	Handle(message *chat.ClientMessage) error
}
