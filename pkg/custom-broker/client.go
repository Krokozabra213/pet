package custombroker

import brokerutils "github.com/Krokozabra213/sso/pkg/custom-broker/utils"

type Client struct {
	uuid uint64
	name string
	buf  chan interface{}
	done chan struct{}
}

type IClient interface {
	GetUUID() uint64
	GetName() string
	GetBuffer() chan interface{}
	GetDone() chan struct{}
	close()
	send(message interface{})
}

func NewClient(name string, bufferSize int) *Client {
	return &Client{
		uuid: brokerutils.GenerateRandomUint64(),
		name: name,
		buf:  make(chan interface{}, bufferSize),
		done: make(chan struct{}),
	}
}

func (cli *Client) GetUUID() uint64 {
	return cli.uuid
}

func (cli *Client) GetName() string {
	return cli.name
}

func (cli *Client) GetBuffer() chan interface{} {
	return cli.buf
}

func (cli *Client) GetDone() chan struct{} {
	return cli.done
}

func (cli *Client) close() {
	if cli == nil {
		return
	}
	close(cli.done)
	close(cli.buf)
}

func (cli *Client) send(message interface{}) {
	if cli == nil {
		return
	}
	select {
	case <-cli.done:
	default:
		cli.buf <- message
	}
}
