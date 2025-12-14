package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	"github.com/Krokozabra213/sso/internal/auth/grpc"
)

func (s *ServerAPI) IsAdmin(
	ctx context.Context,
	r *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {

	isAdminInput := businessinput.NewIsAdminInput(r.GetUserId())
	isAdminOutput, err := s.Business.IsAdmin(ctx, isAdminInput)

	if err != nil {
		return nil, grpc.HandleError(err)
	}

	response := &sso.IsAdminResponse{
		IsAdmin: isAdminOutput.GetAccess(),
	}

	return response, err
}
