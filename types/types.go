package types

import (
	"reflect"

	"go.trulyao.dev/mirror/parser"
)

type ParserInterface interface {
	AddSource(reflect.Type) error
	AddSources(...reflect.Type) error
	ParseN(int) (parser.Item, error)
	LookupByName(string) (parser.Item, bool)

	Next() (parser.Item, error)
	Done() bool
}

// A general language interface to make it harder to pass in a wrong language or extend the built-in languages and backends in the future
// There will clearly be neglibile performance impact but it should not matter much here
type TargetInterface interface {
	Name() string
	Path() string
	Language() string
	Extension() string
	Header() string
	AddCustomType(string, string)
	Generator() GeneratorInterface
}

type GeneratorInterface interface {
	// Assign a parser to use
	SetParser(ParserInterface) error

	// Set a custom header text to prepend at the top of all files
	SetHeaderText(string)

	// Generates the code for the nth element in the parsed items list
	GenerateN(int) (string, error)

	// Generate all types in one go
	GenerateAll() ([]string, error)

	// Generate a single item
	GenerateItem(parser.Item) (string, error)
}

var _ ParserInterface = &parser.Parser{}
