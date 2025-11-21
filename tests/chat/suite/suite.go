package suite

import (
	"context"
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatSuite struct {
	*testing.T
	Conn   *grpc.ClientConn
	Stream chat.Chat_ChatStreamClient
}

func New(t *testing.T) (context.Context, *ChatSuite) {
	t.Helper()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	conn, err := grpc.NewClient(
		"localhost:44045",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	chatClient := chat.NewChatClient(conn)
	stream, err := chatClient.ChatStream(ctx)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	return ctx, &ChatSuite{
		T:      t,
		Stream: stream,
	}
}

func (s *ChatSuite) Close() {
	if s.Conn != nil {
		s.Conn.Close()
	}
}
