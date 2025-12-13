package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	businessoutput "github.com/Krokozabra213/sso/internal/auth/domain/business-output"
)

type IBusiness interface {
	Login(ctx context.Context, input *businessinput.LoginInput) (
		*businessoutput.LoginOutput, error,
	)
	RegisterNewUser(ctx context.Context, input *businessinput.RegisterInput) (
		*businessoutput.RegisterOutput, error,
	)
	PublicKey(ctx context.Context, input *businessinput.PublicKeyInput) (
		*businessoutput.PublicKeyOutput, error,
	)
	IsAdmin(ctx context.Context, input *businessinput.IsAdminInput) (
		*businessoutput.IsAdminOutput, error,
	)
	Logout(ctx context.Context, input *businessinput.LogoutInput) (
		*businessoutput.LogoutOutput, error,
	)
	Refresh(ctx context.Context, input *businessinput.RefreshInput) (
		*businessoutput.RefreshOutput, error,
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

	regInput := businessinput.NewRegisterInput(r.GetUsername(), r.GetPassword())
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

	loginInput := businessinput.NewLoginInput(r.GetUsername(), r.GetPassword(), int(r.GetAppId()))
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

	logoutInput := businessinput.NewLogoutInput(r.GetRefreshToken())
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

	refreshInput := businessinput.NewRefreshInput(r.GetRefreshToken())
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

	isAdminInput := businessinput.NewIsAdminInput(r.GetUserId())
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

	publicKeyInput := businessinput.NewPublicKeyInput(int(r.GetAppId()))
	publicKeyOutput, err := s.Business.PublicKey(ctx, publicKeyInput)

	if err == nil {
		response = &sso.PublicKeyResponse{
			PublicKey: publicKeyOutput.GetPublicKey(),
		}
	}
	return response, err
}
