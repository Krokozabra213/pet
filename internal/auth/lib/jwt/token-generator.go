package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
)

var (
	ErrJWT        = errors.New("failed to generate token")
	ErrParseJWT   = errors.New("failed to parse token")
	ErrUserID     = errors.New("user_id not found in token")
	ErrUsername   = errors.New("username not found in token")
	ErrExp        = errors.New("exp not found in token")
	ErrAppID      = errors.New("app_id not found in token")
	ErrTokenID    = errors.New("err create token_id")
	ErrValidToken = errors.New("err valid token")
)

// Validation claims errors
var (
	ErrValidExp = errors.New("err validation expired")
)

const (
	JWTID     = "jwt_id"
	UserID    = "user_id"
	Username  = "username"
	ExpiredAt = "exp"
	AppID     = "app_id"
)

type IKeyManager interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
}

type IPublicKey interface {
	GetPublicKey() *rsa.PublicKey
}

type TokenGenerator struct {
	user       *domain.User
	app        *domain.App
	accessTTL  time.Duration
	refreshTTL time.Duration
	keyManager IKeyManager
}

func New(
	user *domain.User, app *domain.App,
	accessTTL time.Duration, refreshTTL time.Duration,
	keyManager IKeyManager,
) *TokenGenerator {

	return &TokenGenerator{
		user:       user,
		app:        app,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		keyManager: keyManager,
	}
}

func (gen *TokenGenerator) GenerateTokenPair() (*domain.TokenPair, error) {

	accessToken, err := gen.NewAccess()
	if err != nil {
		return nil, err
	}

	refreshToken, err := gen.NewRefresh()
	if err != nil {
		return nil, err
	}
	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
