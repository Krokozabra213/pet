package chatdomain

import (
	"github.com/Krokozabra213/protos/gen/go/chat"
)

type IConvertServerMessage interface {
	ConvertToServerMessage() *chat.ServerMessage
}

// type IServerMessage interface {
// 	GetUserID() int64
// 	GetUsername() string
// 	GetTimestamp() time.Time
// }
