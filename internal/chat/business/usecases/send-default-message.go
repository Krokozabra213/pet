package chatusecases

import (
	"context"
	"log/slog"

	"github.com/Krokozabra213/sso/internal/chat/domain"
)

func (a *Chat) SendDefaultMessage(ctx context.Context, msg *domain.DefaultMessage) error {
	const op = "chat.SendDefaultMessage-Business"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("send default message", "msg", msg)

	savedMsg, err := a.defaultMsgSaver.SaveDefaultMessage(ctx, msg)
	if err != nil {
		return err
	}

	err = a.msgRepo.Message(ctx, savedMsg)
	if err != nil {
		return err
	}

	return nil
}
