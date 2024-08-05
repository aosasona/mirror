package generator

import "go.trulyao.dev/mirror/parser"

// TODO: implement
type GeneratorInterface interface {
	SetParser(parser.ParserInterface) error
}
