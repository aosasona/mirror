package types

import (
	"reflect"

	"go.trulyao.dev/mirror/parser"
)

type ParserInterface interface {
	// Add a source to the parser
	AddSource(reflect.Type) error

	// Add multiple sources to the parser
	AddSources(...reflect.Type) error

	// Parse the nth source in the list
	ParseN(int) (parser.Item, error)

	// Lookup an type/source by name
	LookupByName(string) (parser.Item, bool)

	// Parse the next source in the list and return the parsed item
	// If there are no sources left to parse, an error will be returned and nil for the item; you should use Done() to check if there are sources left before calling this function to avoid this
	// WARNING: this function consumes the sources and can only be called once
	Next() (parser.Item, error)

	// Check if there are any sources left to parse
	Done() bool

	// Iterate over all parsed sources and apply a function to each
	// If the function returns an error, the iteration will stop and return the error
	// Unlike Next(), this function does not consume the sources and can be called multiple times
	Iterate(func(parser.Item) error) error

	// Count the number of sources left to parse
	Count() int
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

	// Sets whether or not to use strict mode - this is enabled by default for all built-in generators
	// This is mostly useful for testing purposes
	SetNonStrict(bool)
}

var _ ParserInterface = &parser.Parser{}
