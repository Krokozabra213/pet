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
