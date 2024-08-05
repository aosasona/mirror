package generator

import (
	"go.trulyao.dev/mirror/generator/typescript"
	"go.trulyao.dev/mirror/parser"
)

// TODO: implement
type GeneratorInterface interface {
	// Assign a parser to use
	SetParser(parser.ParserInterface) error

	// Set a custom header text to prepend at the top of all files
	SetHeaderText(string)

	// Generates the code for the nth element in the parsed items list
	GenerateN(int) (string, error)

	// Generate all types in one go
	GenerateAll() ([]string, error)
}

func NewTypescriptGenerator() *typescript.Generator {
	return typescript.NewGenerator()
}
