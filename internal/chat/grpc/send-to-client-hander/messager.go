package sendtoclienthander

import (
	"context"
	"log/slog"

	// chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Messager struct {
	log      *slog.Logger
	ctx      context.Context
	buffer   <-chan interface{}
	shutdown <-chan struct{}
}

func NewMessager(
	log *slog.Logger,
	ctx context.Context,
	buffer <-chan interface{},
	shutdown <-chan struct{},
) *Messager {
	return &Messager{
		log:      log,
		ctx:      ctx,
		buffer:   buffer,
		shutdown: shutdown,
	}
}

func (cm *Messager) GetClientMessage() (interface{}, error) {

	select {
	case <-cm.shutdown:
		// Канал закрыт - нормальное завершение
		cm.log.Debug(ErrGracefulShutdown.Error())
		return nil, status.Error(codes.Canceled, ErrGracefulShutdown.Error())
	case message, ok := <-cm.buffer:
		if !ok {
			// Канал закрыт - нормальное завершение
			cm.log.Debug(ErrGracefulShutdown.Error())
			return nil, status.Error(codes.Canceled, ErrGracefulShutdown.Error())
		}
		return message, nil
	case <-cm.ctx.Done():
		// Клиент отключился
		cm.log.Debug(cm.ctx.Err().Error())
		err := HandleStreamContextError(cm.ctx)
		return nil, status.Error(codes.Canceled, err.Error())
	}
}
