package broker

type Client struct {
	ID   uint64
	Name string
	Buf  chan interface{}
	Done chan struct{}
}

func NewClient(id uint64, name string, bufferSize int) Client {
	return Client{
		ID:   id,
		Name: name,
		Buf:  make(chan interface{}, bufferSize),
		Done: make(chan struct{}),
	}
}

func (cli *Client) GetID() uint64 {
	return cli.ID
}

func (cli *Client) GetName() string {
	return cli.Name
}

func (cli *Client) GetBuffer() chan interface{} {
	return cli.Buf
}

func (cli *Client) GetDone() chan struct{} {
	return cli.Done
}

func (cli *Client) Close() {
	close(cli.Done)
	close(cli.Buf)
}

func (cli *Client) Send(message interface{}) {
	select {
	case <-cli.Done:
	default:
		cli.Buf <- message
	}
}

type IClient interface {
	GetID() int64
	GetName() string
	GetBuffer() chan interface{}
	GetDone() chan struct{}
	Close()
	Send(message interface{})
}
