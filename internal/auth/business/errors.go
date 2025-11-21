package authBusiness

import "errors"

var (
	ErrTimeout            = errors.New("timeout request")
	ErrExists             = errors.New("exists error")
	ErrNotFound           = errors.New("not found error")
	ErrInternal           = errors.New("internal service error")
	ErrHashPassword       = errors.New("password hashing error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGenerate      = errors.New("generation access and refresh token error")
	ErrParse              = errors.New("parsing error")
	ErrPermission         = errors.New("not enough permissions")
	ErrTokenExpired       = errors.New("the token has expired")
	ErrTokenRevoked       = errors.New("token revoked")
)
