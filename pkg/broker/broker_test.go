package broker

// import (
// 	"context"
// 	"math/rand"
// 	"strconv"
// 	"testing"
// 	"time"
// )

// type Message struct {
// 	content  string
// 	username string
// }

// func TestBrokerFunction(t *testing.T) {
// 	broker, _ := NewCBroker(3)

// 	ctx := context.Background()

// 	// создаём 100 клиентов, которые будут заходить раз в 0.5сек
// 	clientName := "client"
// 	clients := make([]IClient, 100)
// 	for i := 0; i < 100; i++ {
// 		clients[i] = NewClient(uint64(i), clientName+strconv.Itoa(i), 50)
// 	}

// 	// создаём 1000 сообщений
// 	rand.Seed(time.Now().UnixNano())
// 	cont := "content"
// 	messages := make([]*Message, 1000)
// 	for i := 0; i < 1000; i++ {
// 		messages[i] = &Message{content: cont + strconv.Itoa(i), username: clientName + strconv.Itoa(rand.Intn(100))}
// 	}

// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			broker.Subscribe(clients[i])
// 			time.Sleep(100 * time.Millisecond)
// 		}

// 		for i := 0; i < 100; i++ {
// 			broker.Unsubscribe(clients[i])
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			for i := 0; i < 1000; i++ {
// 				broker.Send(ctx, messages[i])
// 				time.Sleep(10 * time.Millisecond)
// 			}
// 		}
// 	}()

// }

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Message struct {
	content  string
	username string
}

// TestClient расширяет Client для сбора статистики
type TestClient struct {
	*Client
	mu    sync.Mutex
	count int
}

func NewTestClient(id uint64, name string, bufSize int) *TestClient {
	client := &TestClient{
		Client: NewClient(id, name, bufSize),
	}
	return client
}

func (client *TestClient) send(message interface{}) {
	client.mu.Lock()
	client.count++
	client.mu.Unlock()
}

// func TestBrokerFunction(t *testing.T) {
// 	// Статистика
// 	var (
// 		totalSent     int
// 		totalReceived int
// 		subCount      int
// 		unsubCount    int
// 		errorsCount   int
// 		statsMu       sync.Mutex // ОДИН мьютекс для всей статистики
// 	)

// 	broker, _ := NewCBroker(3)
// 	ctx := context.Background()

// 	// Создаём 100 клиентов
// 	clientName := "client"
// 	clients := make([]*TestClient, 100)
// 	for i := 0; i < 100; i++ {
// 		clients[i] = NewTestClient(uint64(i), clientName+strconv.Itoa(i), 100)
// 	}

// 	// Создаём 1000 сообщений
// 	rand.Seed(time.Now().UnixNano())
// 	cont := "content"
// 	messages := make([]*Message, 1000)
// 	for i := 0; i < 1000; i++ {
// 		messages[i] = &Message{
// 			content:  cont + strconv.Itoa(i),
// 			username: clientName + strconv.Itoa(rand.Intn(100)),
// 		}
// 	}

// 	// Канал для завершения теста
// 	done := make(chan bool)

// 	// Таймер для ограничения времени теста
// 	testTimer := time.NewTimer(30 * time.Second)
// 	defer testTimer.Stop()

// 	// Тикер для периодического логирования
// 	statsTicker := time.NewTicker(5 * time.Second)
// 	defer statsTicker.Stop()

// 	// Горутина для логирования статистики
// 	go func() {
// 		for {
// 			select {
// 			case <-statsTicker.C:
// 				statsMu.Lock()
// 				t.Logf("[STATS] Sent: %d, Received: %d, Subscribed: %d, Unsubscribed: %d, Errors: %d",
// 					totalSent, totalReceived, subCount, unsubCount, errorsCount)
// 				statsMu.Unlock()

// 				// Логирование состояния бакетов
// 				for i, bucket := range broker.buckets {
// 					bucket.muDeque.Lock()
// 					dequeSize := len(bucket.deque)
// 					bucket.muDeque.Unlock()

// 					bucket.mu.RLock()
// 					cacheSize := len(bucket.cache)
// 					bucket.mu.RUnlock()

// 					t.Logf("[BUCKET %d] Cache: %d, Deque: %d, ToSendBuf: %d/%d",
// 						i, cacheSize, dequeSize, len(bucket.toSendBuf), cap(bucket.toSendBuf))
// 				}
// 			case <-done:
// 				return
// 			}
// 		}
// 	}()

// 	// Горутина для подписки/отписки клиентов
// 	go func() {
// 		t.Log("Starting subscription cycle...")
// 		for i := 0; i < 100; i++ {
// 			broker.Subscribe(clients[i])
// 			statsMu.Lock()
// 			subCount++
// 			statsMu.Unlock()
// 			t.Logf("Subscribed client %d", i)
// 			time.Sleep(100 * time.Millisecond)
// 		}

// 		time.Sleep(2 * time.Second) // Даем время на отправку сообщений

// 		t.Log("Starting unsubscription cycle...")
// 		for i := 0; i < 100; i++ {
// 			broker.Unsubscribe(clients[i])
// 			statsMu.Lock()
// 			unsubCount++
// 			statsMu.Unlock()
// 			t.Logf("Unsubscribed client %d", i)
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}()

// 	// Горутина для отправки сообщений
// 	go func() {
// 		t.Log("Starting message sending...")
// 		messageIndex := 0
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			default:
// 				if messageIndex >= 1000 {
// 					messageIndex = 0 // циклически отправляем сообщения
// 				}

// 				ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
// 				err := broker.Send(ctxWithTimeout, messages[messageIndex])
// 				cancel()

// 				if err != nil {
// 					statsMu.Lock()
// 					errorsCount++
// 					statsMu.Unlock()
// 					t.Logf("Error sending message %d: %v", messageIndex, err)
// 				} else {
// 					statsMu.Lock()
// 					totalSent++
// 					statsMu.Unlock()
// 					t.Log("message sent...")
// 				}

// 				messageIndex++
// 				time.Sleep(10 * time.Millisecond)
// 			}
// 		}
// 	}()

// 	// Ожидаем завершения теста по таймеру
// 	<-testTimer.C
// 	close(done)

// 	// Даем время на завершение операций
// 	time.Sleep(1 * time.Second)

// 	// Финальная статистика
// 	t.Log("=== FINAL STATISTICS ===")
// 	t.Logf("Total messages sent: %d", totalSent)
// 	t.Logf("Total errors: %d", errorsCount)

// 	// Статистика по клиентам
// 	for i, client := range clients {
// 		count := len(client.messages)
// 		t.Logf("Client %d received %d messages", i, count)
// 	}

// 	// Проверка состояния бакетов
// 	t.Log("=== BUCKETS STATE ===")
// 	totalDequeSize := 0
// 	totalCacheSize := 0
// 	for i, bucket := range broker.buckets {
// 		bucket.muDeque.Lock()
// 		dequeSize := len(bucket.deque)
// 		bucket.muDeque.Unlock()

// 		bucket.mu.RLock()
// 		cacheSize := len(bucket.cache)
// 		bucket.mu.RUnlock()

// 		totalDequeSize += dequeSize
// 		totalCacheSize += cacheSize

// 		if dequeSize > 0 || cacheSize > 0 {
// 			t.Logf("Bucket %d: cache=%d, deque=%d", i, cacheSize, dequeSize)
// 		}
// 	}
// 	t.Logf("Total in all buckets: cache=%d, deque=%d", totalCacheSize, totalDequeSize)
// }

// Дополнительный тест для проверки отправки под большей нагрузкой
func TestBrokerLoad(t *testing.T) {
	broker, _ := NewCBroker(2, 2000) // Меньше бакетов для более концентрированной нагрузки
	ctx := context.Background()

	var successCount int32
	var errorCount int32

	// Создаем клиентов и сразу подписываем
	clients := make([]*TestClient, 500)
	for i := 0; i < 500; i++ {
		clients[i] = NewTestClient(uint64(i), "loadclient"+strconv.Itoa(i), 100)
		broker.subscribe(clients[i])
	}

	time.Sleep(100 * time.Millisecond)

	// Отправляем сообщения параллельно
	var wg sync.WaitGroup
	startTime := time.Now()
	for i := 0; i < 50000; i++ {
		wg.Add(1)
		go func(msgNum int) {
			defer wg.Done()

			msg := &Message{
				content:  fmt.Sprintf("load_message_%d", msgNum),
				username: fmt.Sprintf("user_%d", rand.Intn(50)),
			}

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := broker.send(ctx, msg); err != nil {
				atomic.AddInt32(&errorCount, 1)
			} else {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}
	wg.Wait()
	time.Sleep(5 * time.Second)

	duration := time.Since(startTime)

	t.Logf("Load test completed in %v", duration)
	t.Logf("Successful sends: %d, Errors: %d", successCount, errorCount)
	// t.Logf("Throughput: %.2f messages/second", float64(1000)/duration.Seconds())

	// Проверяем доставку
	totalReceived := 0
	for i, client := range clients {
		totalReceived += client.count

		t.Logf("Client %d received %d messages", i, client.count)
	}

	t.Logf("Total messages received by clients: %d", totalReceived)

	//-----////

	t.Log("=== QUICK DIAGNOSIS ===")

	// Проверяем необработанные сообщения в бакетах
	totalStuck := 0
	for i, bucket := range broker.buckets {
		dequeSize := len(bucket.queue.queue)

		stuckInToSend := len(bucket.toSendBuf)
		stuckInMessages := len(bucket.messages)

		totalStuck += dequeSize + stuckInToSend + stuckInMessages

		if dequeSize > 0 || stuckInToSend > 0 || stuckInMessages > 0 {
			t.Logf("Bucket %d has stuck: deque=%d, toSend=%d, messages=%d",
				i, dequeSize, stuckInToSend, stuckInMessages)
		}
	}

	if totalStuck > 0 {
		t.Logf("Total stuck messages in system: %d", totalStuck)
	}

	t.Log("=== CHECKING MESSAGE STATUSES ===")
	for i, bucket := range broker.buckets {

		statusCount := make(map[uint8]int)
		for _, msg := range bucket.queue.queue {
			// statusCount[(*msg.Status)]++
			statusCount[msg.status]++
		}

		t.Logf("Bucket %d - Status distribution: %v", i, statusCount)
	}

}
