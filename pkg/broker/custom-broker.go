package broker

import (
	"fmt"
	"time"

	"github.com/Krokozabra213/sso/pkg/broker/utils"
)

// maxSpectators - количество горутин, принимающих запросы от пользователей
// maxShards - количество шардированных кешей
// sizeBuffer - размер буфера shard при инициализации
// const (
// maxSpectators = 4
// maxShards     = 5
// )

// BucketsLog - степень двойки количества корзин
func NewCBroker(BucketsLog uint8, msgBufSize, allocationCache, shardSpectatorsCount int) (*CBroker, error) {
	if BucketsLog == 0 {
		return nil, fmt.Errorf("some error")
	}

	// var b uint8
	// if utils.IsPowerOfTwo(countShards) {
	// 	b = countShards
	// } else {
	// 	b = utils.LogarithmFloor(countShards)
	// }

	broker := &CBroker{
		buckets: make([]*Bucket, utils.PowInt(2, BucketsLog)),
		seed:    uint64(time.Now().UnixNano()),
		B:       BucketsLog,
	}
	for i := range len(broker.buckets) {
		broker.buckets[i] = NewBucket(msgBufSize, allocationCache, shardSpectatorsCount)
	}

	return broker, nil
}

type CBroker struct {
	buckets []*Bucket
	seed    uint64
	B       uint8
}

func (CB *CBroker) Subscribe(cli IClient) error {
	hash := utils.SimpleUint64Hash(uint64(cli.GetID()), CB.seed)
	numBuckets := uint64(1) << CB.B
	bucketMask := numBuckets - 1
	bucketIndex := (hash >> (64 - CB.B)) & bucketMask
	err := CB.buckets[bucketIndex].Register(cli)
	// indShard := utils.Hash(cli.GetName(), len(CB.shards))
	// err := CB.shards[indShard].register(cli)
	return err
	// return nil
}

func (CB *CBroker) Send(message interface{}) {
	for _, s := range CB.buckets {
		s.Send(message)
	}
}

func (CB *CBroker) Unsubscribe(cli IClient) {
	// indShard := utils.Hash(cli.GetName(), len(CB.buckets))
	// CB.buckets[indShard].delete(cli)
	hash := utils.SimpleUint64Hash(uint64(cli.GetID()), CB.seed)
	numBuckets := uint64(1) << CB.B
	bucketMask := numBuckets - 1
	bucketIndex := (hash >> (64 - CB.B)) & bucketMask
	CB.buckets[bucketIndex].Delete(cli)
}

// func (CB *CBroker) CheckOnline(cli *ClientDeps) {
// 	CB.once.Do(CB.spectate)
// 	indShard := pkg.Hash(cli.Name, maxShards)
// 	CB.shards[indShard].CheckOnline(cli)
// }
