package chatinterfaces

import (
	"context"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
)

type IChatClient interface {
	GetUUID() uint64
	GetName() string
	GetBuffer() chan interface{}
	GetDone() chan struct{}
}

type IBusiness interface {
	Subscribe(ctx context.Context, username string) (IChatClient, error)
	Unsubscribe(userID uint64)
	ISendImageMessage
	ISendTextMessage
}

type ISendImageMessage interface {
	SendImageMessage(ctx context.Context, msg *chatdomain.ImageMessage) error
}

type ISendTextMessage interface {
	SendTextMessage(ctx context.Context, msg *chatdomain.TextMessage) error
}
