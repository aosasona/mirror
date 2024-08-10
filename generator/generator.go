package generator

import (
	"errors"
)

var (
	ErrNoParser      = errors.New("no parser provided")
	ErrNoName        = errors.New("item has no name")
	ErrUnhandledItem = errors.New("unhandled item type")
	ErrUnknwonType   = errors.New("unknown item type")
)
