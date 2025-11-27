package broker

import (
	"context"
	"sync"
	"time"

	"github.com/Krokozabra213/sso/pkg/broker/utils"
)

const (
	CtxSendTimeout = 3 * time.Second
)

// BucketsLog - степень двойки количества корзин
func NewCBroker(bucketsLog, maxClientCount int) (*CBroker, error) {
	if bucketsLog == 0 || bucketsLog > 16 {
		return nil, ErrBucketsCount
	}

	numBuckets := 1 << bucketsLog // 2^BucketsLog
	maskBuckets := uint64(numBuckets - 1)

	broker := &CBroker{
		buckets: make([]*Bucket, numBuckets),
		seed:    uint64(time.Now().UnixNano()),
		mask:    maskBuckets,
	}

	if numBuckets >= 8 {
		broker.parallelInitBuckets(maxClientCount)
	} else {
		broker.initBuckets(maxClientCount)
	}

	return broker, nil
}

type CBroker struct {
	buckets []*Bucket
	seed    uint64
	mask    uint64
}

func (CB *CBroker) parallelInitBuckets(maxClientCount int) {
	var wg sync.WaitGroup
	for i := 0; i < len(CB.buckets); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			CB.buckets[idx] = NewBucket(maxClientCount)
		}(i)
	}
	wg.Wait()
}

func (CB *CBroker) initBuckets(maxClientCount int) {
	for i := 0; i < len(CB.buckets); i++ {
		CB.buckets[i] = NewBucket(maxClientCount)
	}
}

func (CB *CBroker) subscribe(cli IClient) {
	ind := CB.getBucketIndex(cli.GetUUID())
	CB.buckets[ind].Register(cli)
}

func (CB *CBroker) unsubscribe(cli IClient) {
	ind := CB.getBucketIndex(cli.GetUUID())
	CB.buckets[ind].Delete(cli)
}

func (CB *CBroker) send(ctx context.Context, message interface{}) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, CtxSendTimeout)
		defer cancel()
	}

	isolatingMsg := make([]*wrappedMessage, len(CB.buckets))

	id := utils.GenerateRandomUint64()
	for i := range len(isolatingMsg) {
		isolatingMsg[i] = &wrappedMessage{
			id:      id,
			status:  NotDef, // статус не определен
			message: message,
		}
	}

	var wg sync.WaitGroup
	for i, bucket := range CB.buckets {
		wg.Add(1)
		go func(idx int, b *Bucket) {
			defer wg.Done()
			select {
			case b.messages <- isolatingMsg[idx]:
			case <-ctx.Done():
			}
		}(i, bucket)
	}
	wg.Wait()

	if ctx.Err() != nil {
		for _, m := range isolatingMsg {
			m.status = Cancelled // отменено
		}
		return ctx.Err()
	}

	for _, m := range isolatingMsg {
		m.status = Ready // готово к отправке клиентам
	}

	return nil
}

func (CB *CBroker) getBucketIndex(uuid uint64) uint64 {
	hash := uuid*0x9e3779b97f4a7c15 ^ CB.seed
	return hash & CB.mask
}
