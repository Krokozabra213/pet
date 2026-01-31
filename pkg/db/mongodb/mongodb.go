package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const defaultTimeout = 10 * time.Second

// NewClient создаёт и подключает клиент MongoDB
func NewClient(uri, username, password string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	opts := options.Client().ApplyURI(uri)

	if username != "" && password != "" {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	// mongo.Connect сразу создаёт клиент и инициирует подключение
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Проверяем, что подключение действительно работает
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		// Если ping не прошёл, отключаемся
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	return client, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

// Disconnect корректно закрывает соединение
func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return client.Disconnect(ctx)
}
