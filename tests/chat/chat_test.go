package chat_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/tests/chat/suite"
)

func TestRegisterLogin_Errors(t *testing.T) {
	_, st := suite.New(t)

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := st.Stream.Recv()

			if err == io.EOF {
				t.Log("Server closed the stream")
				close(waitc)
				return
			}
			if err != nil {
				t.Errorf("Failed to receive: %v", err)
				close(waitc)
				return
			}

			// Обрабатываем разные типы сообщений от сервера
			switch msg := in.Type.(type) {
			case *chat.ServerMessage_Joined:
				fmt.Printf("🟢 USER JOINED: %s (ID: %d)\n",
					msg.Joined.Username, msg.Joined.UserId)

			case *chat.ServerMessage_SendMessage:
				fmt.Printf("💬 CHAT MESSAGE: %s (ID: %d) at %s: %s\n",
					msg.SendMessage.Username,
					msg.SendMessage.UserId,
					time.Unix(msg.SendMessage.Timestamp, 0).Format("15:04:05"),
					msg.SendMessage.Content)

			case *chat.ServerMessage_Left:
				fmt.Printf("🔴 USER LEFT: %s (ID: %d)\n",
					msg.Left.Username, msg.Left.UserId)

			default:
				fmt.Printf("❓ UNKNOWN MESSAGE TYPE: %s\n", in.String())
			}
			fmt.Println("==========================")
		}
	}()

	// Отправка сообщений
	messages := []*chat.ClientMessage{
		{ // Первое сообщение - присоединение пользователя
			Type: &chat.ClientMessage_Join{
				Join: &chat.JoinMessage{
					UserId:   1,
					Username: "anton1",
				},
			},
		},
		{ // Второе сообщение - отправка текста
			Type: &chat.ClientMessage_SendMessage{
				SendMessage: &chat.SendMessageAction{
					Content: "content1",
				},
			},
		},
		{ // Третье сообщение - выход пользователя
			Type: &chat.ClientMessage_Leave{
				Leave: &chat.LeaveChat{
					UserId:   1,
					Username: "anton1",
				},
			},
		},
	}

	for _, msg := range messages {
		if err := st.Stream.Send(msg); err != nil {
			t.Fatalf("Failed to send: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	st.Stream.CloseSend()
	<-waitc // Ждем завершения чтения
}
