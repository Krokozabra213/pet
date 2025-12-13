package chatusecases

import (
	"context"
	"log/slog"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
)

func (a *Chat) SendImageMessage(ctx context.Context, msg *chatdomain.ImageMessage) error {
	const op = "chat.SendImageMessage-Business"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("send image message", "msg", msg)

	var err error
	msg, err = a.msgSaver.SaveImageMessage(ctx, msg)
	if err != nil {
		return err
	}

	err = a.msgRepo.Message(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
