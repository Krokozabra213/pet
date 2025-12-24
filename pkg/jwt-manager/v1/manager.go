package jwtv1

import (
	"crypto/rsa"
	"time"
)

// access&refresh claims
const (
	JWTID     = "jwt_id"
	UserID    = "user_id"
	Username  = "username"
	ExpiredAt = "exp"
	AppID     = "app_id"
)

type Data struct {
	UserID   uint64
	Username string
	AppID    uint
}

type JWTManager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New(public *rsa.PublicKey, private *rsa.PrivateKey) (*JWTManager, error) {
	if public == nil {
		return nil, ErrEmptyPublicKey
	}

	if private == nil {
		return nil, ErrEmptyPrivateKey
	}

	return &JWTManager{publicKey: public, privateKey: private}, nil
}

func (m *JWTManager) GenerateTokens(
	accessTTL, refreshTTL time.Duration, data *Data,
) (string, string, error) {

	accessToken, err := m.GenerateAccess(
		accessTTL, data.UserID, data.Username, data.AppID,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := m.GenerateRefresh(
		refreshTTL, data.UserID, data.Username, data.AppID,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
