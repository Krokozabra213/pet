package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"google.golang.org/grpc"
)

const (
	emptyVal = 0
)

type IBusiness interface {
	Login(ctx context.Context, username, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, username, password string) (userID int64, err error)
	AppSault(ctx context.Context, appID int) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	Business IBusiness
}

func Register(grpc *grpc.Server, business IBusiness) {
	sso.RegisterAuthServer(grpc, &serverAPI{Business: business})
}

func (s *serverAPI) Register(
	ctx context.Context,
	r *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	username, pass := r.Username, r.Password
	if username == "" || pass == "" {
		return nil, ErrInvalidCredentials
	}

	userID, err := s.Business.RegisterNewUser(ctx, username, pass)
	return &sso.RegisterResponse{
		UserId: userID,
	}, err
}

func (s *serverAPI) Login(
	ctx context.Context,
	r *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	username, pass, appID := r.Username, r.Password, r.AppId
	if username == "" || pass == "" || appID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	token, err := s.Business.Login(ctx, username, pass, int(appID))
	return &sso.LoginResponse{
		Token: token,
	}, err
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	r *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	userID := r.UserId
	if userID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	access, err := s.Business.IsAdmin(ctx, userID)
	return &sso.IsAdminResponse{
		IsAdmin: access,
	}, err
}

func (s *serverAPI) AppSault(
	ctx context.Context,
	r *sso.AppSaultRequest,
) (*sso.AppSaultResponse, error) {

	appID := r.AppId
	if appID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	sault, err := s.Business.AppSault(ctx, int(appID))
	return &sso.AppSaultResponse{
		Sault: sault,
	}, err
}
