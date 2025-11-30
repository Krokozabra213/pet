package chatgrpc

import "errors"

// server
var (
	ErrStream       = errors.New("failed to receive message")
	ErrDisconect    = errors.New("client disconnected")
	ErrFirstMessage = errors.New("first message must be join")
)

// utils
var (
	ErrSendMessage        = errors.New("error send message")
	ErrRecvMessage        = errors.New("error recv message")
	ErrUnknownMessageType = errors.New("unknown message type")
)
