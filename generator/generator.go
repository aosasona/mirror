package generator

import (
	"errors"
)

var (
	ErrNoFields      = errors.New("no fields present in struct")
	ErrNoParser      = errors.New("no parser provided")
	ErrNoName        = errors.New("item has no name")
	ErrNoBaseItem    = errors.New("item has no base item")
	ErrUnhandledItem = errors.New("unhandled item type")
	ErrUnknownType   = errors.New("unknown item type")
)
