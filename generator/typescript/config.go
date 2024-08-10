package typescript

import (
	"go.trulyao.dev/mirror/config"
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

	// AllowUnexportedFields will include private fields
	AllowUnexportedFields bool

	// IndentationType is the type of indentation to use (space or tab)
	IndentationType config.Indentation

	// IndentationCount is the number of spaces or tabs to use for indentation (defaults to 4)
	IndentationCount int

	customTypes map[string]string
}

func New(filename, path string) *Config {
	return &Config{
		FileName:         filename,
		OutputPath:       path,
		customTypes:      make(map[string]string),
		IndentationCount: 4,
	}
}

func (c *Config) Name() string { return c.FileName }

func (c *Config) Path() string { return c.OutputPath }

func (c *Config) Language() string { return "typescript" }

func (c *Config) Extension() string { return "ts" }

func (c *Config) Header() string {
	return fileHeader
}

func (c *Config) SetFileName(name string) *Config {
	c.FileName = name
	return c
}

func (c *Config) SetOutputPath(path string) *Config {
	c.OutputPath = path
	return c
}

func (c *Config) SetPreferNullForNullable(value bool) *Config {
	c.PreferNullForNullable = value
	return c
}

func (c *Config) SetInlineObjects(value bool) *Config {
	c.InlineObjects = value
	return c
}

func (c *Config) SetIncludeSemiColon(value bool) *Config {
	c.InludeSemiColon = value
	return c
}

func (c *Config) SetPreferUnknown(value bool) *Config {
	c.PreferUnknown = value
	return c
}

func (c *Config) SetAllowUnexportedFields(value bool) *Config {
	c.AllowUnexportedFields = value
	return c
}

func (c *Config) SetIndentationType(value config.Indentation) *Config {
	c.IndentationType = value
	return c
}

func (c *Config) SetIndentationCount(value int) *Config {
	c.IndentationCount = value
	return c
}

func (c *Config) AddCustomType(name, value string) {
	c.customTypes[name] = value
}

func (c *Config) Generator() *Generator {
	return NewGenerator(c)
}
