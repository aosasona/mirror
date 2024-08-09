package typescript

import (
	"go.trulyao.dev/mirror/generator"
	"go.trulyao.dev/mirror/types"
)

type Generator struct {
	config Config
	parser types.ParserInterface
}

// SetParser implements generator.GeneratorInterface.
func (g *Generator) SetParser(parser types.ParserInterface) error {
	if parser == nil {
		return generator.ErrNoParser
	}

	g.parser = parser
	return nil
}

// GenerateAll implements generator.GeneratorInterface.
func (g *Generator) GenerateAll() ([]string, error) {
	panic("unimplemented")
}

// GenerateN implements generator.GeneratorInterface.
func (g *Generator) GenerateN(int) (string, error) {
	panic("unimplemented")
}

// SetHeaderText implements generator.GeneratorInterface.
func (g *Generator) SetHeaderText(string) {
	panic("unimplemented")
}

func NewGenerator(config Config) *Generator {
	return &Generator{config: config}
}
