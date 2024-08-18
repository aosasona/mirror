package types

import (
	"reflect"

	"go.trulyao.dev/mirror/v2/parser"
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

	// Set the parser's configuration
	SetConfig(parser.Config) error

	// Reset the parser to its initial state
	Reset()
}

// A general language interface to make it harder to pass in a wrong language or extend the built-in languages and backends in the future
// There will clearly be neglibile performance impact but it should not matter much here
type TargetInterface interface {
	// Returns the target file name with the extension
	Name() string

	// Returns the target file path
	Path() string

	// Returns the target language name (e.g. "typescript")
	Language() string

	// Returns the target file extension (e.g. ".ts")
	Extension() string

	// Returns the target file header text
	Header() string

	// Prefix for the types in the target
	Prefix() string

	// Add a custom type to the target
	AddCustomType(string, string)

	// Create and return a new instance of the language's generator based on the target's config
	Generator() GeneratorInterface

	// Unique identifier for the target
	ID() string

	// Check if two targets are equivalent (i.e. have the same ID)
	IsEquivalent(TargetInterface) bool

	// Ensure the target has been configured correctly
	Validate() error
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
