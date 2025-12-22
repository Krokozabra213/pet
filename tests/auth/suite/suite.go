package suite

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	ssonewconfig "github.com/Krokozabra213/sso/newconfigs/sso"
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
	Cfg        *ssonewconfig.Config
	AuthClient sso.AuthClient
	DB         *postgrespet.PGDB
	Redis      *redispet.RDB
}

func New(t *testing.T) (context.Context, *SSOSuite) {
	t.Helper()

	cfg, err := ssonewconfig.Init("settings/sso_main.yml", "sso.env")
	cfg.PG.DSN = "host=0.0.0.0 user=user password=password dbname=postgres port=5555 sslmode=disable"
	cfg.Redis.Addr = "0.0.0.0:6379"
	if err != nil {
		t.Fatalf("config init err: %v", err)
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	DB := postgrespet.NewPGDB(cfg.PG.DSN)
	redis := redispet.NewRedisDB(cfg.Redis.Addr, cfg.Redis.Pass, cfg.Redis.Cache)

	cc, err := grpc.NewClient(
		grpcAddress(cfg.GRPC.Host, cfg.GRPC.Port),
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
	result := s.DB.Exec("TRUNCATE TABLE users CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SSOSuite) CleanupAppsData() error {
	// удаляем все строки в таблице apps
	result := s.DB.Exec("TRUNCATE TABLE apps CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SSOSuite) CleanupAdminsData() error {
	// удаляем все строки в таблице admins
	result := s.DB.Exec("TRUNCATE TABLE admins CASCADE")
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
	result := s.DB.Create(app)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(app.ID), nil
}

func (s *SSOSuite) CreateAdmin(userID int64) (int64, error) {
	admin := &domain.Admin{
		UserID: userID,
	}
	result := s.DB.Create(admin)
	if result.Error != nil {
		return 0, result.Error
	}
	return int64(admin.UserID), nil
}

func grpcAddress(host, port string) string {
	return net.JoinHostPort(host, port)
}
