package business

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrSchoolNotFound    = errors.New("school not found")
	ErrNoAccess          = errors.New("no access")
)
