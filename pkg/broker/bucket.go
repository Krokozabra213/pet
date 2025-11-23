package broker

import (
	"fmt"
	"sync"
)

type Bucket struct {
	mu         sync.RWMutex
	cache      Users
	messages   chan interface{}
	spectators int
}

func NewBucket(msgBufSize, allocationCache, shardSpectatorsCount int) *Bucket {
	bucket := &Bucket{
		cache:      make(Users, allocationCache),
		messages:   make(chan interface{}, msgBufSize),
		spectators: shardSpectatorsCount,
	}
	bucket.broadcast()
	return bucket
}

func (S *Bucket) Register(cli IClient) error {
	S.mu.Lock()
	defer S.mu.Unlock()
	if _, exists := S.cache[cli.GetID()]; !exists {
		S.cache[cli.GetID()] = cli
	}
	return fmt.Errorf("пользователь с именем %s уже на сервере", cli.GetName())
}

func (S *Bucket) Send(message interface{}) {
	S.messages <- message
}

func (S *Bucket) Delete(cli IClient) {
	cli.Close()
	delete(S.cache, cli.GetID())
}

func (S *Bucket) broadcast() {
	for range S.spectators {
		go func() {

			for m := range S.messages {
				S.mu.RLock()
				defer S.mu.RUnlock()
				for _, client := range S.cache {
					client.Send(m)
				}
			}
		}()
	}
}
