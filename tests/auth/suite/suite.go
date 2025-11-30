package suite

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

type SSOSuite struct {
	*testing.T
	Cfg        *ssoconfig.Config
	AuthClient sso.AuthClient
	DB         *postgrespet.PGDB
	Redis      *redispet.RDB
}

func New(t *testing.T) (context.Context, *SSOSuite) {
	t.Helper()

	env := EnvLocal

	cfg := ssoconfig.Load(env, true)
	t.Logf("Config: %+v", cfg)

	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	DB := postgrespet.NewPGDB(cfg.DB.DSN)
	redis := redispet.NewRedisDB(cfg.Redis.Addr, cfg.Redis.Pass, cfg.Redis.Cache)

	cc, err := grpc.NewClient(
		grpcAddress(cfg.Server.Host, cfg.Server.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatalf("grpc server connection failed %v", err)
	}

	return ctx, &SSOSuite{
		T:          t,
		Cfg:        cfg,
		AuthClient: sso.NewAuthClient(cc),
		DB:         DB,
		Redis:      redis,
	}
}

func (s *SSOSuite) CleanupTestData() error {
	err := s.CleanupUserData()
	if err != nil {
		return err
	}
	err = s.CleanupBlackTokensData()
	if err != nil {
		return err
	}
	err = s.CleanupAppsData()
	if err != nil {
		return err
	}
	err = s.CleanupAdminsData()
	if err != nil {
		return err
	}
	err = s.CleanupRedis()
	if err != nil {
		return err
	}
	return nil
}

func (s *SSOSuite) CleanupUserData() error {
	// удаляем все строки в таблице users
	result := s.DB.Client.Exec("TRUNCATE TABLE users CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// поменять на redis
func (s *SSOSuite) CleanupBlackTokensData() error {
	// удаляем все строки в таблице black_tokens
	result := s.DB.Client.Exec("TRUNCATE TABLE black_tokens CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SSOSuite) CleanupAppsData() error {
	// удаляем все строки в таблице apps
	result := s.DB.Client.Exec("TRUNCATE TABLE apps CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SSOSuite) CleanupAdminsData() error {
	// удаляем все строки в таблице admins
	result := s.DB.Client.Exec("TRUNCATE TABLE admins CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *SSOSuite) CleanupRedis() error {
	ctx := context.Background()
	err := r.Redis.Client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush database: %w", err)
	}
	return nil
}

func (s *SSOSuite) CreateApp(name string) (int, error) {
	app := &domain.App{
		Name: name,
	}
	result := s.DB.Client.Create(app)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(app.ID), nil
}

func (s *SSOSuite) CreateAdmin(userID int64) (int64, error) {
	admin := &domain.Admin{
		UserID: userID,
	}
	result := s.DB.Client.Create(admin)
	if result.Error != nil {
		return 0, result.Error
	}
	return int64(admin.UserID), nil
}

func grpcAddress(host, port string) string {
	return net.JoinHostPort(host, port)
}
