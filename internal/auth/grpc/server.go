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

	// TODO: VALIDATE
	username, pass := r.Username, r.Password
	if username == "" || pass == "" {
		return nil, ErrInvalidCredentials
	}

	userID, err := s.Business.RegisterNewUser(ctx, username, pass)
	return &sso.RegisterResponse{
		UserId: int64(userID),
	}, err
}

func (s *serverAPI) Login(
	ctx context.Context,
	r *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	// TODO: VALIDATE
	username, pass, appID := r.Username, r.Password, r.AppId
	if username == "" || pass == "" || appID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.Business.Login(ctx, username, pass, int(appID))
	return &sso.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

func (s *serverAPI) Logout(
	ctx context.Context,
	r *sso.LogoutRequest,
) (*sso.LogoutResponse, error) {

	// TODO: VALIDATE
	refreshToken := r.RefreshToken
	if refreshToken == "" {
		return nil, ErrInvalidCredentials
	}

	success, err := s.Business.Logout(ctx, refreshToken)

	return &sso.LogoutResponse{
		Success: success,
	}, err
}

func (s *serverAPI) Refresh(
	ctx context.Context,
	r *sso.RefreshRequest,
) (*sso.RefreshResponse, error) {

	// TODO: VALIDATE
	token := r.RefreshToken
	if token == "" {
		return nil, ErrInvalidCredentials
	}

	access, refresh, err := s.Business.Refresh(ctx, token)
	return &sso.RefreshResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, err
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	r *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	// TODO: VALIDATE
	userID := r.UserId
	if userID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	access, err := s.Business.IsAdmin(ctx, userID)
	return &sso.IsAdminResponse{
		IsAdmin: access,
	}, err
}

func (s *serverAPI) GetPublicKey(
	ctx context.Context,
	r *sso.PublicKeyRequest,
) (*sso.PublicKeyResponse, error) {

	// TODO: VALIDATE
	appID := r.AppId
	if appID == emptyVal {
		return nil, ErrInvalidCredentials
	}

	publicKey, err := s.Business.PublicKey(ctx, int(appID))
	return &sso.PublicKeyResponse{
		PublicKey: publicKey,
	}, err
}
