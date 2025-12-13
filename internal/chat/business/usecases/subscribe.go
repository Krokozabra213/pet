package chatusecases

import (
	"context"
	"log/slog"

	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"

	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

func (a *Chat) Subscribe(ctx context.Context, username string) (chatgrpc.IChatClient, error) {

	const op = "chat.Subscribe-Business"
	log := a.log.With(
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
