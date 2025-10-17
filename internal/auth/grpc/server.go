package grpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"google.golang.org/grpc"
)

const (
	emptyVal = 0
)

// TODO

// type Business interface {
// 	Login(ctx context.Context, username, password string, appID int) (token string, err error)
// 	Register(ctx context.Context, username, password string) (userID int64, err error)
// 	GetAppSault(ctx context.Context, appID int) (string, error)
// 	IsAdmin(ctx context.Context, userID int64) (bool, error)
// }

type serverAPI struct {
	sso.UnimplementedAuthServer
}

func Register(grpc *grpc.Server) {
	sso.RegisterAuthServer(grpc, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	r *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	username, pass, appID := r.Username, r.Password, r.AppId
	if username == "" || pass == "" || appID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	resp := &sso.LoginResponse{
		Token: "test",
	}
	return resp, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	r *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	username, pass := r.Username, r.Password
	if username == "" || pass == "" {
		return nil, ErrInvalidCredentials
	}

	resp := &sso.RegisterResponse{
		UserId: 1,
	}
	return resp, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	r *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	userID := r.UserId
	if userID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	resp := &sso.IsAdminResponse{
		IsAdmin: true,
	}
	return resp, nil
}

func (s *serverAPI) AppSault(
	ctx context.Context,
	r *sso.AppSaultRequest,
) (*sso.AppSaultResponse, error) {

	appID := r.AppId
	if appID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	resp := &sso.AppSaultResponse{
		Sault: "test",
	}
	return resp, nil
}
