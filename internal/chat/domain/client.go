package domain

type Client struct {
	ID   int64
	Name string
	Buf  chan interface{}
	Done chan struct{}
}
