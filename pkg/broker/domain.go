package broker

import (
	"log"
	"sync"
)

// message statuses
const (
	NotDef    uint8 = 0 // 0 - не определен
	Ready     uint8 = 1 // 1 - готово к отправке
	Cancelled uint8 = 2 // 2 - отменено
	Sended    uint8 = 3 // 3 - отправлено
)

const (
	LenCleanUp = 400
	BatchSize  = 250
)

type wrappedMessage struct {
	id      uint64
	status  uint8
	message interface{}
}

type usersMap map[uint64]IClient
type queueMap map[uint64]*wrappedMessage

// memory cache пользователей чата
type clientCache struct {
	mu      sync.Mutex
	clients usersMap
}

func newClientCache(cap int) *clientCache {
	return &clientCache{
		clients: make(usersMap, cap),
	}
}

func (cc *clientCache) register(cli IClient) {
	cc.mu.Lock()
	cc.clients[cli.GetUUID()] = cli
	cc.mu.Unlock()
}

func (cc *clientCache) delete(cli IClient) {
	cc.mu.Lock()
	delete(cc.clients, cli.GetUUID())
	cc.mu.Unlock()
}

func (cc *clientCache) distributingMessages(message interface{}) {
	cc.mu.Lock()
	for _, client := range cc.clients {
		client.send(message)
	}
	cc.mu.Unlock()
}

// memory cache сообщений чата
type messageCache struct {
	mu        sync.Mutex
	queue     queueMap
	queueSize int
	toSendBuf chan interface{}
}

func newMessageCache(cap int, toSendBuf chan interface{}) *messageCache {
	return &messageCache{
		queue:     make(queueMap, cap),
		queueSize: cap,
		toSendBuf: toSendBuf,
	}
}

func (mc *messageCache) cleanupDeque() {

	if len(mc.queue) < LenCleanUp {
		return
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	newQueue := make(queueMap, mc.queueSize)
	for k, v := range mc.queue {
		if v.status < Cancelled { // Сохраняем только неотправленные и неопределенные
			newQueue[k] = v
		}
	}
	mc.queue = newQueue
}

func (mc *messageCache) set(m *wrappedMessage) {
	mc.mu.Lock()
	if len(mc.queue) > mc.queueSize {
		log.Printf("WARNING: large deque size: %d", len(mc.queue))
	}
	mc.queue[m.id] = m
	mc.mu.Unlock()
}

func (mc *messageCache) processDequeBatch() {
	mc.mu.Lock()
	for _, v := range mc.queue {
		if v.status == Ready {
			v.status = Sended
			select {
			case mc.toSendBuf <- v.message:
				// успех
			default:
				v.status = Ready
			}
		}
	}
	mc.mu.Unlock()
}
