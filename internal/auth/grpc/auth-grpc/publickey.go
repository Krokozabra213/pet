package authgrpc

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/sso"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	"github.com/Krokozabra213/sso/internal/auth/grpc"
)

func (s *ServerAPI) GetPublicKey(
	ctx context.Context,
	r *sso.PublicKeyRequest,
) (*sso.PublicKeyResponse, error) {

	publicKeyInput := businessinput.NewPublicKeyInput(int(r.GetAppId()))
	publicKeyOutput, err := s.Business.PublicKey(ctx, publicKeyInput)

	if err != nil {
		return nil, grpc.HandleError(err)
	}

	response := &sso.PublicKeyResponse{
		PublicKey: publicKeyOutput.GetPublicKey(),
	}

	return response, err
}
