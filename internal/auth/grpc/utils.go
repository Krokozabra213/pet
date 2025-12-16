package grpc

import (
	"errors"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	var bizErr authBusiness.IBusinessError
	if errors.As(err, &bizErr) {
		return bizErr.ToGRPC()
	}

	return status.Error(codes.Internal, "internal error")
}
