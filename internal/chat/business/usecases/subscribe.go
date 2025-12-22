package chatusecases

import (
	"context"
	"log/slog"

	chatinterfaces "github.com/Krokozabra213/sso/internal/chat/grpc/interfaces"

	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

func (a *Chat) Subscribe(ctx context.Context, username string) (chatinterfaces.IChatClient, error) {

	const op = "chat.Subscribe-Business"
	log := slog.With(
		slog.String("op", op),
	)

	client := custombroker.NewClient(username, BufferSize)
	err := a.clientRepo.Subscribe(ctx, client)
	if err != nil {
		log.Error("failed subscribe", "err", err)
		return nil, err
	}
	return client, nil
}
