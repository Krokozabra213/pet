package authBusiness

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user already exists")
	ErrUnknown            = errors.New("unknown error")
	ErrPermission         = errors.New("permission error")
)
