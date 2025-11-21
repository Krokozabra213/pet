package broker

import (
	"fmt"
	"sync"
)

type Shard struct {
	mu        sync.RWMutex
	cache     Users
	messages  chan string
	online    chan *ClientDeps
	addClient chan *ClientDeps
	delete    chan *ClientDeps
}

func NewShard(sizeBuffer int) *Shard {
	return &Shard{
		cache:     make(Users, sizeBuffer),
		messages:  make(chan string, sizeBuffer),
		online:    make(chan *ClientDeps, sizeBuffer),
		addClient: make(chan *ClientDeps, sizeBuffer),
		delete:    make(chan *ClientDeps, sizeBuffer),
	}
}

func (S *Shard) Register(cli *ClientDeps) error {
	S.mu.RLock()
	defer S.mu.RUnlock()
	if _, exist := S.cache[cli.Name]; exist {
		return fmt.Errorf("пользователь с именем %s уже на сервере", cli.Name)
	}
	S.addClient <- cli
	return nil
}

func (S *Shard) Send(message string) {
	S.messages <- message
}

func (S *Shard) CheckOnline(cli *ClientDeps) {
	S.online <- cli
}

func (S *Shard) Delete(cli *ClientDeps) {
	S.delete <- cli
}

func (S *Shard) broadcast() {
	for range maxSpectators {
		go func() {
			for {
				select {
				case client := <-S.addClient:
					S.mu.Lock()
					S.cache[client.Name] = client.Cli
					S.mu.Unlock()

				case client := <-S.delete:
					S.mu.Lock()
					delete(S.cache, client.Name)
					S.mu.Unlock()

				case client := <-S.online:
					// не синхронизируем, т.к. не значимая операция
					if cli, exist := S.cache[client.Name]; exist {
						go cli.sendMessage(S.cache.Online())
					}

				case m := <-S.messages:
					S.mu.RLock()
					cache := S.cache
					S.mu.RUnlock()

					for _, user := range cache {
						user.sendMessage(m)
					}
				}
			}
		}()
	}
}
