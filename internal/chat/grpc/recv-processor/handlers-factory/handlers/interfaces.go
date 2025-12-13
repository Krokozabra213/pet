package handlers

import "github.com/Krokozabra213/protos/gen/go/chat"

type IMessageHandler interface {
	GetType() string
	CanHandle(message *chat.ClientMessage) bool
	Handle(message *chat.ClientMessage) error
}
