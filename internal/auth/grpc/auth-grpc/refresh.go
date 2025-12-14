package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	"github.com/Krokozabra213/sso/internal/auth/grpc"
)

func (s *ServerAPI) Refresh(
	ctx context.Context,
	r *sso.RefreshRequest,
) (*sso.RefreshResponse, error) {
	refreshInput := businessinput.NewRefreshInput(r.GetRefreshToken())
	refreshOutput, err := s.Business.Refresh(ctx, refreshInput)

	if err != nil {
		return nil, grpc.HandleError(err)
	}

	response := &sso.RefreshResponse{
		AccessToken:  refreshOutput.GetAccess(),
		RefreshToken: refreshInput.GetRefreshToken(),
	}

	return response, err
}
