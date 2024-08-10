package typescript

import (
	"go.trulyao.dev/mirror/generator"
	"go.trulyao.dev/mirror/parser"
	"go.trulyao.dev/mirror/types"
)

type Generator struct {
	config *Config
	parser types.ParserInterface
}

func NewGenerator(config *Config) *Generator {
	return &Generator{config: config}
}

// SetParser implements generator.GeneratorInterface.
func (g *Generator) SetParser(parser types.ParserInterface) error {
	if parser == nil {
		return generator.ErrNoParser
	}

	g.parser = parser
	return nil
}

// GenerateItem implements types.GeneratorInterface.
func (g *Generator) GenerateItem(item parser.Item) (string, error) {
	switch item := item.(type) {
	case parser.Scalar:
		return g.generateScalar(item)
	case parser.List:
	// TODO: implement
	case parser.Map:
	// TODO: implement
	case parser.Struct:
	// TODO: implement
	case parser.Function:
		return g.generateFunction(item)
	default:
		return "", generator.ErrUnknwonType
	}

	return "", generator.ErrUnhandledItem
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

func (g *Generator) generateScalar(item parser.Scalar) (string, error) {
	if item.Name == "" {
		return "", generator.ErrNoName
	}

	var (
		typeString = "type %s = %s"
		typeName   = item.Name
		typeValue  string
	)

	switch item.Type {
	case parser.TypeAny:
		typeValue = "any"

	case parser.TypeString:
		typeValue = "string"

	case parser.TypeInteger, parser.TypeFloat:
		typeValue = "number"
	}
}

func (g *Generator) generateStruct(item parser.Struct) (string, error) {
	panic("unimplemented")
}

func (g *Generator) generateList(item parser.List) (string, error) {
	panic("unimplemented")
}

func (g *Generator) generateArray(item parser.List) (string, error) {
	panic("unimplemented")
}

func (g *Generator) generateMap(item parser.List) (string, error) {
	panic("unimplemented")
}

func (g *Generator) generateFunction(item parser.Function) (string, error) {
	panic("unimplemented")
}

var _ types.GeneratorInterface = &Generator{}
