package pgxclient

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host              string
	Port              string
	User              string
	Password          string
	Database          string
	SSLMode           string
	ConnectTimeout    time.Duration
	MaxConns          int32
	MinConns          int32
	MaxConnLifeTime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

// Значения по умолчанию
func NewDefaultConfig() Config {
	return Config{
		Host:              "localhost",
		Port:              "5432",
		SSLMode:           "disable",
		ConnectTimeout:    5 * time.Second,
		MaxConns:          10,
		MinConns:          2,
		MaxConnLifeTime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

func (c Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.User == "" {
		return fmt.Errorf("user is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database is required")
	}
	if c.MaxConns < 1 {
		return fmt.Errorf("max_conns must be >= 1")
	}
	if c.MinConns > c.MaxConns {
		return fmt.Errorf("min_conns cannot exceed max_conns")
	}
	return nil
}

// DSN с экранированием пароля
func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

type Client struct {
	pool *pgxpool.Pool
	cfg  Config
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	// Валидация
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifeTime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Таймаут на подключение
	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Проверяем соединение
	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Client{
		pool: pool,
		cfg:  cfg,
	}, nil
}

// NewWithRetry — подключение с повторными попытками
func NewWithRetry(ctx context.Context, cfg Config, maxRetries int, retryDelay time.Duration) (*Client, error) {
	var client *Client
	var err error

	for i := 0; i <= maxRetries; i++ {
		client, err = New(ctx, cfg)
		if err == nil {
			return client, nil
		}

		if i < maxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay):
				// продолжаем
			}
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

// Shutdown — graceful shutdown с ожиданием завершения запросов
func (c *Client) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		c.pool.Close()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout: %w", ctx.Err())
	case <-done:
		return nil
	}
}

// Close — немедленное закрытие (для defer)
func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *Client) Pool() *pgxpool.Pool {
	return c.pool
}

// Stats — статистика пула для мониторинга
func (c *Client) Stats() *pgxpool.Stat {
	return c.pool.Stat()
}

// Health — проверка здоровья для health checks
func (c *Client) Health(ctx context.Context) error {
	if err := c.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database unhealthy: %w", err)
	}

	stat := c.pool.Stat()
	if stat.TotalConns() == 0 {
		return fmt.Errorf("no available connections")
	}

	return nil
}
