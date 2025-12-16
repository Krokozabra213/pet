package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	"github.com/Krokozabra213/sso/internal/auth/grpc"
)

func (s *ServerAPI) Logout(
	ctx context.Context,
	r *sso.LogoutRequest,
) (*sso.LogoutResponse, error) {

	logoutInput := businessinput.NewLogoutInput(r.GetRefreshToken())
	logoutOutput, err := s.Business.Logout(ctx, logoutInput)

	if err != nil {
		return nil, grpc.HandleError(err)
	}

	response := &sso.LogoutResponse{
		Success: logoutOutput.GetSuccess(),
	}

	return response, err
}
