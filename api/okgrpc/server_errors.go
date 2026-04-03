package okgrpc

import "errors"

var (
	ErrInvalidPayload    = errors.New("invalid payload")
	ErrInvalidUpdateMask = errors.New("invalid update mask")
)
