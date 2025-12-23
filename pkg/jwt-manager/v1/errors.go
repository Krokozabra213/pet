package jwtv1

import "errors"

var (
	ErrEmptyPublicKey     = errors.New("empty public key")
	ErrEmptyPrivateKey    = errors.New("empty private key")
	ErrSignedAccessToken  = errors.New("failed sign access jwt token")
	ErrSignedRefreshToken = errors.New("failed sign refresh jwt token")
	ErrGenerateJWTID      = errors.New("failed generate jwt_id")
)

// var (
// 	ErrJWT        = errors.New("failed to generate token")
// 	ErrParseJWT   = errors.New("failed to parse token")
// 	ErrUserID     = errors.New("user_id not found in token")
// 	ErrUsername   = errors.New("username not found in token")
// 	ErrExp        = errors.New("exp not found in token")
// 	ErrAppID      = errors.New("app_id not found in token")
// 	ErrTokenID    = errors.New("err create token_id")
// 	ErrValidToken = errors.New("err valid token")
// )

// var (
// 	ErrValidExp = errors.New("err validation expired")
// )
