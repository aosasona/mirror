package target

import "go.trulyao.dev/mirror/parser"

// TODO: implement
type TargetInterface interface {
	SetParser(parser.Parser) error
}
