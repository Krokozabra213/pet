package custombroker

import (
	"context"
	"sync"
	"time"

	brokerutils "github.com/Krokozabra213/sso/pkg/custom-broker/utils"
	WP "github.com/Krokozabra213/sso/pkg/custom-broker/wrapped-message"
	// "github.com/Krokozabra213/sso/pkg/broker/utils"
	// WP "github.com/Krokozabra213/sso/pkg/broker/wrapped-message"
)

const (
	CtxSendTimeout            = 3 * time.Second
	GetCurrentClientCountTick = 1 * time.Second
)

type CBroker struct {
	buckets        []*Bucket
	seed           uint64
	mask           uint64
	maxClientCount int
	clientCount    int
}

// bucketsLog - логарифмический показатель количества корзин (bucket'ов) брокера.
// Определяет количество корзин по формуле: numBuckets = 2^bucketsLog
//
// Примеры:
//   bucketsLog = 1 -> 2 корзины
//   bucketsLog = 4 -> 16 корзин
//   bucketsLog = 8 -> 256 корзин

// Технические детали:
//   - Используется для распределения клиентов по корзинам по хешу: bucket = hash(client) & mask
//   - Маска вычисляется как: mask = (1 << bucketsLog) - 1
//   - Такой подход позволяет использовать быстрые битовые операции вместо деления

// maxClientCount - максимальное количество клиентов, которые смогут одновременно пользоваться чатом
// Допускается погрешность, при большом количестве запросов Subscribe, т.к. проверка на количество клиентов
// производится каждые 1 сек (const GetCurrentClientCountTick)

// multcachesmemmory - множитель памяти, для алокации структур, отвечающих за хранение сообщений и клиентов
// Считается по формуле: (maxClientCount/numBuckets)*multcachesmemmory

func NewCBroker(bucketsLog, maxClientCount, multcachesmemmory int) (*CBroker, error) {
	if bucketsLog == 0 || bucketsLog > 16 {
		return nil, ErrBucketsCount
	}

	numBuckets := 1 << bucketsLog // 2^BucketsLog
	maskBuckets := uint64(numBuckets - 1)

	broker := &CBroker{
		buckets:        make([]*Bucket, numBuckets),
		seed:           uint64(time.Now().UnixNano()), // Для рандомизации UUID для Message
		mask:           maskBuckets,
		maxClientCount: maxClientCount,
	}

	// Заранее подготовленная память для структур данных
	cap := int(float64(maxClientCount/numBuckets) * float64(multcachesmemmory))

	if numBuckets >= 8 {
		broker.parallelInitBuckets(cap)
	} else {
		broker.initBuckets(cap)
	}

	// Периодичная одновременная проверка корзин на количество клиентов
	go broker.startCurrentClientWorker()

	return broker, nil
}

func (CB *CBroker) parallelInitBuckets(cap int) {
	var wg sync.WaitGroup
	for i := 0; i < len(CB.buckets); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			CB.buckets[idx] = NewBucket(cap)
		}(i)
	}
	wg.Wait()
}

func (CB *CBroker) initBuckets(cap int) {
	for i := 0; i < len(CB.buckets); i++ {
		CB.buckets[i] = NewBucket(cap)
	}
}

func (CB *CBroker) Subscribe(cli IClient) error {
	if CB.clientCount > CB.maxClientCount {
		return ErrServerIsFull
	}
	ind := CB.getBucketIndex(cli.GetUUID())
	CB.buckets[ind].register(cli)
	return nil
}

func (CB *CBroker) Unsubscribe(cli IClient) {
	ind := CB.getBucketIndex(cli.GetUUID())
	CB.buckets[ind].delete(cli)
}

// Получаем ctx, message от клиента -> отправляем во все buckets ->
// -> проверяем ctx.Err -> меняем статус сообщения
func (CB *CBroker) Send(ctx context.Context, message interface{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		// Добавляем deadline контексту, если отсутствует
		ctx, cancel = context.WithTimeout(ctx, CtxSendTimeout)
		defer cancel()
	}

	// для доступа к полям структуры по ссылке
	isolatingMsg := make([]*WP.WrappedMessage, len(CB.buckets))

	id := brokerutils.GenerateRandomUint64()
	for i := range len(isolatingMsg) {
		isolatingMsg[i] = WP.New(id, NotDef, message) // NotDef - статус не определен
	}

	var wg sync.WaitGroup
	for i, bucket := range CB.buckets {
		wg.Add(1)
		go func(idx int, b *Bucket) {
			defer wg.Done()
			select {
			case b.inputMessageChan <- isolatingMsg[idx]:
			case <-ctx.Done():
			}
		}(i, bucket)
	}
	wg.Wait()

	if ctx.Err() != nil {
		for _, m := range isolatingMsg {
			m.Status = Cancelled // отменено
		}
		return ctx.Err()
	}

	for _, m := range isolatingMsg {
		m.Status = Ready // готово к отправке клиентам
	}

	return nil
}

func (CB *CBroker) getBucketIndex(uuid uint64) uint64 {
	hash := uuid*0x9e3779b97f4a7c15 ^ CB.seed
	return hash & CB.mask // 1100 & 1010 = 1000 (8 в десятичной)
}

func (CB *CBroker) startCurrentClientWorker() {
	ticker := time.NewTicker(GetCurrentClientCountTick)
	for range ticker.C {
		CB.getCurrentClientCount()
	}
}

func (CB *CBroker) getCurrentClientCount() {
	count := 0
	ch := make(chan int, len(CB.buckets))

	var wg sync.WaitGroup
	for _, b := range CB.buckets {
		wg.Add(1)
		go func(bucket *Bucket) {
			defer wg.Done()
			ch <- bucket.getClientCount()
		}(b)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for c := range ch {
		count += c
	}

	CB.clientCount = count
}
