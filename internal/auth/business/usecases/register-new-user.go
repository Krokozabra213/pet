package authusecases

import (
	"context"
	"errors"
	"log/slog"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	businessoutput "github.com/Krokozabra213/sso/internal/auth/domain/business-output"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"golang.org/x/crypto/bcrypt"
)

func (a *Auth) RegisterNewUser(
	ctx context.Context, input *businessinput.RegisterInput,
) (*businessoutput.RegisterOutput, error) {
	const op = "auth.RegisterNewUser"

	username := input.GetUsername()
	password := input.GetPassword()

	log := slog.With(
		slog.String("op", op),
		slog.String("username", username),
	)
	log.Info("starting user registration process")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed password hashing", slog.String("error", err.Error()))
		return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrHashPassword)
	}

	user := domain.NewUser(username, string(passHash))
	userID, err := a.userProvider.SaveUser(ctx, user)
	if err != nil {
		log.Error("failed save new user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrDuplicate) {
			return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrExists)
		}
		return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrInternal)
	}

	log.Info("user successfully registered",
		slog.Uint64("user_id", uint64(userID)))

	uotput := businessoutput.NewRegisterOutput(userID)
	return uotput, nil
}
