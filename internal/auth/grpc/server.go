package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"github.com/Krokozabra213/sso/internal/auth/domain"
)

type IBusiness interface {
	// Login(ctx context.Context, username, password string, appID int) (
	// 	accessToken string, refreshToken string, err error,
	// )
	// RegisterNewUser(ctx context.Context, username, password string) (
	// 	userID uint, err error,
	// )
	// PublicKey(ctx context.Context, appID int) (string, error)
	// IsAdmin(ctx context.Context, userID int64) (bool, error)
	// Logout(ctx context.Context, refreshToken string) (bool, error)
	// Refresh(ctx context.Context, token string) (
	// 	accessToken string, refreshToken string, err error,
	// )
	Login(ctx context.Context, input *domain.LoginInput) (
		*domain.LoginOutput, error,
	)
	RegisterNewUser(ctx context.Context, input *domain.RegisterInput) (
		*domain.RegisterOutput, error,
	)
	PublicKey(ctx context.Context, input *domain.PublicKeyInput) (
		*domain.PublicKeyOutput, error,
	)
	IsAdmin(ctx context.Context, input *domain.IsAdminInput) (
		*domain.IsAdminOutput, error,
	)
	Logout(ctx context.Context, input *domain.LogoutInput) (
		*domain.LogoutOutput, error,
	)
	Refresh(ctx context.Context, input *domain.RefreshInput) (
		*domain.RefreshOutput, error,
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
	var response *sso.RegisterResponse

	regInput := domain.NewRegisterInput(r.GetUsername(), r.GetPassword())
	regOutput, err := s.Business.RegisterNewUser(ctx, regInput)

	if err == nil {
		response = &sso.RegisterResponse{
			UserId: int64(regOutput.GetUserID()),
		}
	}
	return response, err
}

func (s *ServerAPI) Login(
	ctx context.Context,
	r *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	var response *sso.LoginResponse

	loginInput := domain.NewLoginInput(r.GetUsername(), r.GetPassword(), int(r.GetAppId()))
	loginOutput, err := s.Business.Login(ctx, loginInput)

	if err == nil {
		response = &sso.LoginResponse{
			AccessToken:  loginOutput.GetAccess(),
			RefreshToken: loginOutput.GetRefresh(),
		}
	}
	return response, err
}

func (s *ServerAPI) Logout(
	ctx context.Context,
	r *sso.LogoutRequest,
) (*sso.LogoutResponse, error) {
	var response *sso.LogoutResponse

	logoutInput := domain.NewLogoutInput(r.GetRefreshToken())
	logoutOutput, err := s.Business.Logout(ctx, logoutInput)

	if err == nil {
		response = &sso.LogoutResponse{
			Success: logoutOutput.GetSuccess(),
		}
	}
	return response, err
}

func (s *ServerAPI) Refresh(
	ctx context.Context,
	r *sso.RefreshRequest,
) (*sso.RefreshResponse, error) {
	var response *sso.RefreshResponse

	refreshInput := domain.NewRefreshInput(r.GetRefreshToken())
	refreshOutput, err := s.Business.Refresh(ctx, refreshInput)

	if err == nil {
		response = &sso.RefreshResponse{
			AccessToken:  refreshOutput.GetAccess(),
			RefreshToken: refreshInput.GetRefreshToken(),
		}
	}
	return response, err
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context,
	r *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	var response *sso.IsAdminResponse

	isAdminInput := domain.NewIsAdminInput(r.GetUserId())
	isAdminOutput, err := s.Business.IsAdmin(ctx, isAdminInput)

	if err == nil {
		response = &sso.IsAdminResponse{
			IsAdmin: isAdminOutput.GetAccess(),
		}
	}
	return response, err
}

func (s *ServerAPI) GetPublicKey(
	ctx context.Context,
	r *sso.PublicKeyRequest,
) (*sso.PublicKeyResponse, error) {
	var response *sso.PublicKeyResponse

	publicKeyInput := domain.NewPublicKeyInput(int(r.GetAppId()))
	publicKeyOutput, err := s.Business.PublicKey(ctx, publicKeyInput)

	if err == nil {
		response = &sso.PublicKeyResponse{
			PublicKey: publicKeyOutput.GetPublicKey(),
		}
	}
	return response, err
}
