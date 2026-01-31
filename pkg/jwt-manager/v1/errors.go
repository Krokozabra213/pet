package jwtv1

import "errors"

var (
	ErrEmptyPublicKey     = errors.New("empty public key")
	ErrEmptyPrivateKey    = errors.New("empty private key")
	ErrSignedAccessToken  = errors.New("failed sign access jwt token")
	ErrSignedRefreshToken = errors.New("failed sign refresh jwt token")
	ErrGenerateJWTID      = errors.New("failed generate jwt_id")
	ErrValidExp           = errors.New("err validation expired")
	ErrParseJWT           = errors.New("failed to parse token")
	ErrValidToken         = errors.New("err valid token")
	ErrUserID             = errors.New("user_id not found in token")
	ErrUsername           = errors.New("username not found in token")
	ErrExp                = errors.New("exp not found in token")
	ErrAppID              = errors.New("app_id not found in token")
	ErrJWTID              = errors.New("jwt_id not found in token")
	ErrInvalidClaims      = errors.New("invalid claims")
)
