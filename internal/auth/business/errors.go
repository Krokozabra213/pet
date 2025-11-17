package authBusiness

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user already exists")
	ErrPermission         = errors.New("permission error")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrHashPassword       = errors.New("hash password error")
	ErrTokenExpired       = errors.New("token expired")
	ErrTimeout            = errors.New("timeout request")
)

//domain errors
var (
	ErrAppUnknown   = errors.New("unknown app error")
	ErrUserUnknown  = errors.New("unknown user error")
	ErrTokenUnknown = errors.New("unknown token error")
)
