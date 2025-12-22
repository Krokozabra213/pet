package chatusecases

import (
	"context"
	"log/slog"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
)

func (a *Chat) SendTextMessage(ctx context.Context, msg *chatdomain.TextMessage) error {
	const op = "chat.SendTextMessage-Business"
	log := slog.With(
		slog.String("op", op),
	)
	log.Info("send text message", "msg", msg)

	var err error
	msg, err = a.msgSaver.SaveTextMessage(ctx, msg)
	if err != nil {
		return err
	}

	err = a.msgRepo.Message(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
