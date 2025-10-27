package configs

import "errors"

var (
	ErrLoadConfig = errors.New("error load config")
	ErrEmptyValue = errors.New("error empty value")
)
