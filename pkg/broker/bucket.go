package broker

import (
	"time"
)

const (
	garbageCollectorTick = 200 * time.Millisecond
	dequeProcessTicker   = 100 * time.Millisecond
)

type Bucket struct {
	cache    *clientCache
	messages chan *wrappedMessage

	queue     *messageCache
	toSendBuf chan interface{}

	// dequeSize int
}

func NewBucket(maxClientCount int) *Bucket {
	cap := float32(maxClientCount) * 0.625
	toSendBuf := make(chan interface{}, int(cap))

	bucket := &Bucket{
		cache:     newClientCache(int(cap)),
		messages:  make(chan *wrappedMessage, int(cap)),
		queue:     newMessageCache(int(cap), toSendBuf),
		toSendBuf: toSendBuf,
	}
	bucket.startWorkers()
	return bucket
}

func (S *Bucket) startWorkers() {
	go S.messageReceiver()
	go S.dequeProcessor()
	go S.messageSender()
	go S.garbageCollector()
}

func (S *Bucket) garbageCollector() {
	ticker := time.NewTicker(garbageCollectorTick)

	for range ticker.C {
		S.cleanupDeque()
	}
}

func (S *Bucket) cleanupDeque() {
	S.queue.cleanupDeque()
}

func (S *Bucket) messageReceiver() {
	for m := range S.messages {
		S.queue.set(m)
	}
}

func (S *Bucket) messageSender() {
	// принимаем сообщения с буфера для сообщения со статусом Ready
	for m := range S.toSendBuf {
		S.cache.distributingMessages(m)
	}
}

func (S *Bucket) dequeProcessor() {
	ticker := time.NewTicker(dequeProcessTicker)
	for range ticker.C {
		S.processDequeBatch()
	}
}

func (S *Bucket) processDequeBatch() {
	S.queue.processDequeBatch()
}

func (S *Bucket) Register(cli IClient) {
	S.cache.register(cli)
}

func (S *Bucket) Delete(cli IClient) {
	S.cache.delete(cli)
}
