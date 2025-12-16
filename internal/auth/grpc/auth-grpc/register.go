package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	"github.com/Krokozabra213/sso/internal/auth/grpc"
)

func (s *ServerAPI) Register(
	ctx context.Context,
	r *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	regInput := businessinput.NewRegisterInput(r.GetUsername(), r.GetPassword())
	regOutput, err := s.Business.RegisterNewUser(ctx, regInput)

	if err != nil {
		return nil, grpc.HandleError(err)
	}

	response := &sso.RegisterResponse{
		UserId: int64(regOutput.GetUserID()),
	}

	return response, err
}
