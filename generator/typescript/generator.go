package typescript

import (
	"go.trulyao.dev/mirror/generator"
	"go.trulyao.dev/mirror/parser"
)

type Generator struct {
	parser parser.ParserInterface
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) SetParser(p parser.ParserInterface) error {
	return nil
}

var _ generator.GeneratorInterface = &Generator{}
