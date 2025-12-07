package custombroker

import (
	"sync"
)

type usersMap map[uint64]IClient

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

func (cc *clientCache) Len() int {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return len(cc.clients)
}

func (cc *clientCache) register(cli IClient) {
	cc.mu.Lock()
	cc.clients[cli.GetUUID()] = cli
	cc.mu.Unlock()
}

func (cc *clientCache) delete(uuid uint64) {
	cc.mu.Lock()
	if client, ok := cc.clients[uuid]; ok {
		client.close()
		delete(cc.clients, uuid)
	}
	cc.mu.Unlock()
}

func (cc *clientCache) distributingMessages(message interface{}) {
	cc.mu.Lock()
	for _, client := range cc.clients {
		client.send(message)
	}
	cc.mu.Unlock()
}
