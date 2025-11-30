package custombroker

import "errors"

var (
	ErrBucketsCount = errors.New("bucketslog must be between 1 and 16")
	ErrServerIsFull = errors.New("server is full")
)
