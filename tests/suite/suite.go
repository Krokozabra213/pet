package suite

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOSuite struct {
	*testing.T
	Cfg        *ssoconfig.Config
	AuthClient sso.AuthClient
	DB         *db.Db
}

func New(t *testing.T) (context.Context, *SSOSuite) {
	t.Helper()
	t.Parallel()

	cfg := ssoconfig.Load("local", true)

	// ctx, cancelCtx := context.WithTimeout(context.Background(), time.Duration(cfg.Server.TimeOut))
	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	db := db.NewPGDb(cfg.DB.DSN)

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
		DB:         db,
	}
}

func (s *SSOSuite) CleanupTestData() error {
	// удаляем все строки в таблице users
	result := s.DB.DB.Exec("TRUNCATE TABLE apps CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SSOSuite) CreateApp(name string) (int, error) {
	app := &domain.App{
		Name: name,
	}
	result := s.DB.DB.Create(app)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(app.ID), nil
}

func grpcAddress(host, port string) string {
	return net.JoinHostPort(host, port)
}
