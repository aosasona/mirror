package typescript

import (
	"errors"
	"path"
	"strings"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/types"
)

type Config struct {
	// FileName is the name of the generated file
	FileName string

	// OutputPath is the path to write the generated file to
	OutputPath string

	// PreferNullForNullable will prefer `null` over `undefined` for nullable types
	PreferNullForNullable bool

	// PreferArrayGeneric will prefer `Array<T>` over `T[]`
	PreferArrayGeneric bool

	// InlineObjects will inline object types instead of using the name (e.g foo: { bar: string } instead of foo: Bar)
	InlineObjects bool

	// InludeSemiColon will include a semi-colon at the end of each type definition
	InludeSemiColon bool

	// PreferUnknown will prefer `unknown` over `any`
	PreferUnknown bool

	// IndentationType is the type of indentation to use (space or tab)
	IndentationType config.Indentation

	// IndentationCount is the number of spaces or tabs to use for indentation (defaults to 4)
	IndentationCount int

	// Prefix is the prefix to add to the generated types (e.g. type Person -> type MyPrefixPerson)
	TypePrefix string

	customTypes map[string]string
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	return &Config{
		FileName:              "generated",
		OutputPath:            "./",
		PreferNullForNullable: true,
		PreferArrayGeneric:    true,
		InlineObjects:         false,
		InludeSemiColon:       true,
		PreferUnknown:         false,
		IndentationType:       config.IndentSpace,
		IndentationCount:      4,
		customTypes:           make(map[string]string),
	}
}

// New returns a new Config with the provided filename and path
func New(filename, path string) *Config {
	return &Config{
		FileName:         filename,
		OutputPath:       path,
		customTypes:      make(map[string]string),
		IndentationCount: 4,
	}
}

// ID returns a unique identifier for a target
func (c *Config) ID() string {
	return strings.ReplaceAll(path.Join(c.OutputPath, c.Name()), "/", ":")
}

// IsEquivalent checks if two targets are equivalent
func (c *Config) IsEquivalent(target types.TargetInterface) bool {
	return c.ID() == target.ID()
}

// Prefix returns the prefix to add to the generated types
func (c *Config) Prefix() string {
	return c.TypePrefix
}

// Name returns the name of the file
func (c *Config) Name() string {
	fileName := c.FileName
	if strings.HasSuffix(fileName, ".ts") {
		return fileName
	}

	return c.FileName + ".ts"
}

// Path returns the path to write the file to
func (c *Config) Path() string {
	return c.OutputPath
}

// Language returns the target language
func (c *Config) Language() string { return "typescript" }

// Extension returns the file extension
func (c *Config) Extension() string { return "ts" }

// Header returns the header text for the file
func (c *Config) Header() string { return fileHeader }

// SetFileName sets the name of the file to write to
func (c *Config) SetFileName(name string) *Config {
	c.FileName = name
	return c
}

// SetOutputPath sets the path to write the file to
func (c *Config) SetOutputPath(path string) *Config {
	c.OutputPath = path
	return c
}

// SetPreferNullForNullable sets whether or not to prefer `null` over `undefined` for nullable types
func (c *Config) SetPreferNullForNullable(value bool) *Config {
	c.PreferNullForNullable = value
	return c
}

// SetInlineObjects sets whether or not to inline object types instead of using the name
// this will result in `foo: { bar: string }` instead of `foo: Bar`
// This is useful for generating types that you do not want to include in the generated file as a separate type
func (c *Config) SetInlineObjects(value bool) *Config {
	c.InlineObjects = value
	return c
}

// SetIncludeSemiColon sets whether or not to include a semi-colon at the end of each type definition
func (c *Config) SetIncludeSemiColon(value bool) *Config {
	c.InludeSemiColon = value
	return c
}

// SetPreferUnknown sets whether or not to prefer `unknown` over `any`
func (c *Config) SetPreferUnknown(value bool) *Config {
	c.PreferUnknown = value
	return c
}

// SetIndentationType sets the type of indentation to use (space or tab)
func (c *Config) SetIndentationType(value config.Indentation) *Config {
	c.IndentationType = value
	return c
}

// SetIndentationCount sets the number of spaces or tabs to use for indentation (defaults to 4)
func (c *Config) SetIndentationCount(value int) *Config {
	c.IndentationCount = value
	return c
}

// SetPrefix sets the prefix to add to the generated types
func (c *Config) SetPrefix(value string) *Config {
	c.TypePrefix = value
	return c
}

// AddCustomType adds a custom type to the config
func (c *Config) AddCustomType(name, value string) {
	c.customTypes[name] = value
}

// Generator returns a new Generator for the current language with the config
func (c *Config) Generator() types.GeneratorInterface {
	return NewGenerator(c)
}

func (c *Config) Validate() error {
	if c.FileName == "" {
		return errors.New("no file name provided")
	}

	if c.OutputPath == "" {
		return errors.New("no output path provided")
	}

	if c.IndentationCount < 2 {
		return errors.New("indentation count must be greater than or equal to 2")
	}

	if c.IndentationType != config.IndentSpace && c.IndentationType != config.IndentTab {
		return errors.New(
			"invalid indentation type, expected `config.IndentSpace` or `config.IndentTab` ",
		)
	}

	return nil
}
