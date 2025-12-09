package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
)

type IBusiness interface {
	Login(ctx context.Context, username, password string, appID int) (
		accessToken string, refreshToken string, err error,
	)
	RegisterNewUser(ctx context.Context, username, password string) (
		userID uint, err error,
	)
	PublicKey(ctx context.Context, appID int) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	Logout(ctx context.Context, refreshToken string) (bool, error)
	Refresh(ctx context.Context, token string) (
		accessToken string, refreshToken string, err error,
	)
}

type ServerAPI struct {
	sso.UnimplementedAuthServer
	Business IBusiness
}

func New(business IBusiness) *ServerAPI {
	return &ServerAPI{
		Business: business,
	}
}

func (s *ServerAPI) Register(
	ctx context.Context,
	r *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	userID, err := s.Business.RegisterNewUser(ctx, r.GetUsername(), r.GetPassword())
	return &sso.RegisterResponse{
		UserId: int64(userID),
	}, err
}

func (s *ServerAPI) Login(
	ctx context.Context,
	r *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	access, refresh, err := s.Business.Login(ctx, r.GetUsername(), r.GetPassword(), int(r.GetAppId()))
	return &sso.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, err
}

func (s *ServerAPI) Logout(
	ctx context.Context,
	r *sso.LogoutRequest,
) (*sso.LogoutResponse, error) {

	success, err := s.Business.Logout(ctx, r.GetRefreshToken())
	return &sso.LogoutResponse{
		Success: success,
	}, err
}

func (s *ServerAPI) Refresh(
	ctx context.Context,
	r *sso.RefreshRequest,
) (*sso.RefreshResponse, error) {

	access, refresh, err := s.Business.Refresh(ctx, r.GetRefreshToken())
	return &sso.RefreshResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, err
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context,
	r *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	access, err := s.Business.IsAdmin(ctx, r.GetUserId())
	return &sso.IsAdminResponse{
		IsAdmin: access,
	}, err
}

func (s *ServerAPI) GetPublicKey(
	ctx context.Context,
	r *sso.PublicKeyRequest,
) (*sso.PublicKeyResponse, error) {

	publicKey, err := s.Business.PublicKey(ctx, int(r.GetAppId()))
	return &sso.PublicKeyResponse{
		PublicKey: publicKey,
	}, err
}
