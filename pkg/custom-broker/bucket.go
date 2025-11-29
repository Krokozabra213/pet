package custombroker

import (
	"log"
	"time"

	WP "github.com/Krokozabra213/sso/pkg/custom-broker/wrapped-message"
)

const (
	GarbageCollectorTick = 200 * time.Millisecond // интервал очистки кеша сообщений от устаревших
	// DequeProcessTicker - интервал отправки сообщений в буфер кеша -> отправка клиентам
	DequeProcessTicker = 100 * time.Millisecond
)

type Bucket struct {
	inputMessageChan chan *WP.WrappedMessage // буфер сообщений broker -> bucket
	clientCache      *clientCache            // in memory cache клиентов
	messageQueue     *messageCache           // in memory cache сообщений
}

func NewBucket(cap int) *Bucket {

	bucket := &Bucket{
		inputMessageChan: make(chan *WP.WrappedMessage, cap),
		clientCache:      newClientCache(cap),
		messageQueue:     newMessageCache(cap),
	}
	bucket.startWorkers()
	return bucket
}

// количество клиентов в кеше
func (S *Bucket) getClientCount() int {
	return S.clientCache.Len()
}

func (S *Bucket) startWorkers() {
	go S.messageReceiver()  // broker:(message)->bucket:inputMessageChan->messagecache map[]
	go S.dequeProcessor()   // messagecache: map[messages...]->chan messages
	go S.messageSender()    // chan messages -> map[clients...]
	go S.garbageCollector() // clean map[messages...]
}

// получаем сообщения broker -> bucket
// добавляем в кеш сообщений
func (S *Bucket) messageReceiver() {
	for m := range S.inputMessageChan {
		S.messageQueue.set(m)
	}
}

// перебираем сообщения кеша (проверяем статус), отправляем в канал
func (S *Bucket) dequeProcessor() {
	ticker := time.NewTicker(DequeProcessTicker)
	for range ticker.C {
		S.messageQueue.processDequeBatch()
	}
}

// рассылаем сообщения клиентам
func (S *Bucket) messageSender() {
	// принимаем сообщения со статусом Ready
	for m := range S.messageQueue.getBuffer() {
		S.clientCache.distributingMessages(m)
	}
}

// сборщик отменённых и отправленных сообщений
func (S *Bucket) garbageCollector() {
	ticker := time.NewTicker(GarbageCollectorTick)

	for range ticker.C {
		S.messageQueue.cleanupDeque()
	}
}

func (S *Bucket) register(cli IClient) {
	log.Printf("user %s joined/n", cli.GetName())
	S.clientCache.register(cli)
}

func (S *Bucket) delete(uuid uint64) {
	S.clientCache.delete(uuid)
}
