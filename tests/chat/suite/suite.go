package suite

import (
	"context"
	"testing"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatSuite struct {
	*testing.T
	Conn    *grpc.ClientConn
	Streams []chat.Chat_ChatStreamClient
}

func New(t *testing.T) (context.Context, *ChatSuite) {
	t.Helper()

	ctx := context.Background()

	conn, err := grpc.NewClient(
		"localhost:44045",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	clients := make([]chat.ChatClient, 0, 50)
	for i := 0; i < 50; i++ {
		clients = append(clients, chat.NewChatClient(conn))
	}

	streams := make([]chat.Chat_ChatStreamClient, 0, 50)
	for i := 0; i < 50; i++ {
		stream, err := clients[i].ChatStream(ctx)
		if err != nil {
			t.Fatalf("Failed to create stream: %v", err)
		}
		streams = append(streams, stream)
	}

	return ctx, &ChatSuite{
		T:       t,
		Streams: streams,
	}
}

func (s *ChatSuite) Close() {
	if s.Conn != nil {
		s.Conn.Close()
	}
}
