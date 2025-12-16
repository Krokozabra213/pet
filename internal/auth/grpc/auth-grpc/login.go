package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	"github.com/Krokozabra213/sso/internal/auth/grpc"
)

func (s *ServerAPI) Login(
	ctx context.Context,
	r *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	loginInput := businessinput.NewLoginInput(r.GetUsername(), r.GetPassword(), int(r.GetAppId()))
	loginOutput, err := s.Business.Login(ctx, loginInput)

	if err != nil {
		return nil, grpc.HandleError(err)
	}

	response := &sso.LoginResponse{
		AccessToken:  loginOutput.GetAccess(),
		RefreshToken: loginOutput.GetRefresh(),
	}
	return response, nil
}
