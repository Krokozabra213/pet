package broker

import (
	"chat-server/pkg/broker/pkg"
	"errors"
	"sync"
)

// maxSpectators - количество горутин, принимающих запросы от пользователей
// maxShards - количество шардированных кешей
// sizeBuffer - размер буфера shard при инициализации
const (
	maxSpectators = 4
	maxShards     = 5
)

func NewCBroker(sizeBuffer int) *CBroker {
	broker := &CBroker{
		shards: make([]*Shard, maxShards),
	}
	for i := range maxShards {
		broker.shards[i] = NewShard(sizeBuffer)
	}
	return broker
}

type CBroker struct {
	once   sync.Once
	shards []*Shard
}

func (CB *CBroker) Register(cli *ClientDeps) error {
	if cli.Name == "" {
		return errors.New("не валидное имя пользователя")
	}
	CB.once.Do(CB.spectate)
	indShard := pkg.Hash(cli.Name, maxShards)
	err := CB.shards[indShard].Register(cli)
	return err
}

func (CB *CBroker) Send(message string) {
	for _, s := range CB.shards {
		s.Send(message)
	}
}

func (CB *CBroker) CheckOnline(cli *ClientDeps) {
	CB.once.Do(CB.spectate)
	indShard := pkg.Hash(cli.Name, maxShards)
	CB.shards[indShard].CheckOnline(cli)
}

func (CB *CBroker) Delete(cli *ClientDeps) {
	CB.once.Do(CB.spectate)
	indShard := pkg.Hash(cli.Name, maxShards)
	CB.shards[indShard].Delete(cli)
}

func (CB *CBroker) Spectate() {
	CB.once.Do(CB.spectate)
}

func (CB *CBroker) spectate() {
	for _, s := range CB.shards {
		s.broadcast()
	}
}
