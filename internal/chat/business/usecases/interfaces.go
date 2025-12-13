package chatusecases

import (
	"context"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

type IClientRepo interface {
	Subscribe(ctx context.Context, client custombroker.IClient) error
	Unsubscribe(uuid uint64)
}

type IMessageRepo interface {
	Message(ctx context.Context, message interface{}) error
}

type IMessageSaver interface {
	ITextMessageSaver
	IImageMessageSaver
}

type ITextMessageSaver interface {
	SaveTextMessage(context.Context, *chatdomain.TextMessage) (*chatdomain.TextMessage, error)
}

type IImageMessageSaver interface {
	SaveImageMessage(context.Context, *chatdomain.ImageMessage) (*chatdomain.ImageMessage, error)
}
