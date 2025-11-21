package broker

import "time"

type Client struct {
	Message chan string
	Out     chan struct{}
}

func NewClient() Client {
	return Client{
		Message: make(chan string),
		Out:     make(chan struct{}),
	}
}

func (cli Client) Close() {
	close(cli.Out)
	// закрываем канал с задержкой чтобы не было ошибки, если канал message закроется раньше чем out
	// и сообщение отправится в закрытый канал
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(cli.Message)
	}()
}

func (cli Client) sendMessage(message string) {

	select {
	case <-cli.Out:
	default:
		cli.Message <- message
	}
}
