package sendtoclienthander

import (
	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertMessage(message interface{}) (*chat.ServerMessage, error) {
	switch msg := message.(type) {

	case domain.IServerMessage:
		convertedMessage := msg.ConvertToServerMessage()
		return convertedMessage, nil

	default:
		return nil, status.Error(codes.InvalidArgument, ErrUnknownMessageType.Error())
	}
}
