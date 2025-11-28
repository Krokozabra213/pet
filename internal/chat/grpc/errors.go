package chatgrpc

import "errors"

var (
	ErrStream       = errors.New("failed to receive message")
	ErrDisconect    = errors.New("client disconnected")
	ErrFirstMessage = errors.New("first message must be join")
)
