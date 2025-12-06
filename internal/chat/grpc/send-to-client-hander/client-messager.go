package sendtoclienthander

import (
	"context"

	// chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Messager struct {
	ctx      context.Context
	buffer   <-chan interface{}
	shutdown <-chan struct{}
}

func NewMessager(
	ctx context.Context,
	buffer <-chan interface{},
	shutdown <-chan struct{},
) *Messager {
	return &Messager{
		ctx:      ctx,
		buffer:   buffer,
		shutdown: shutdown,
	}
}

func (cm *Messager) GetClientMessage() (interface{}, error) {

	select {
	case <-cm.shutdown:
		// Канал закрыт - нормальное завершение
		return nil, status.Error(codes.Internal, ErrGracefulShutdown.Error())
	case message, ok := <-cm.buffer:
		if !ok {
			// Канал закрыт - нормальное завершение
			return nil, status.Error(codes.Internal, ErrGracefulShutdown.Error())
		}
		return message, nil
	case <-cm.ctx.Done():
		// Клиент отключился
		err := HandleStreamContextError(cm.ctx)
		return nil, status.Error(codes.Canceled, err.Error())
	}
}
