package broker

type Broker struct {
	B interface{}
}

func New() *Broker {
	return &Broker{}
}

func (br *Broker) Subscribe(
	id int64, name string, buffer chan interface{},
) error {
	// todo
	return nil
}

func (br *Broker) Unsubscribe(id int64) error {
	// todo
	return nil
}

func (br *Broker) SendMessage(msg string) error {
	// todo
	return nil
}
