package chatgrpc

import (
	"fmt"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) ValidateJoinMessage(message *chat.ClientMessage) (*chat.ClientMessage_Join, error) {
	const op = "chat.Handler.validateJoinMessage"
	log := s.Log.With(
		slog.String("op", op),
	)

	joinMsg, ok := message.Type.(*chat.ClientMessage_Join)
	if !ok {
		log.Error(ErrUnknownMessageType.Error(), slog.String("type", fmt.Sprintf("%T", message.Type)))
		return nil, status.Error(codes.InvalidArgument, ErrFirstMessage.Error())
	}
	return joinMsg, nil
}
