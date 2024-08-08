package config

type Indentation int

const (
	emptyIndent Indentation = iota
	Space
	Tab
)

// A general language interface to make it harder to pass in a wrong language or extend the built-in languages and backends in the future
// There will clearly be neglibile performance impact but it should not matter much here
type TargetInterface interface {
	Name() string
	Path() string
	Language() string
	Extension() string
	Header() string
	AddCustomType(string, string)
	// Generator() generator.GeneratorInterface
}

// Debug is a global variable that can be used to enable or disable debug mode
var Debug = false

func SetDebug(v bool) {
	Debug = v
}

// Pointers have been used here to make sure the user actually sets the value and not just uses the default value
type Config struct {
	// Enabled can be used to disable or enable the generation of types, defaults to false
	Enabled bool

	// Targets are the languages and files to generate types for, at least ONE target MUST be defined
	Targets []TargetInterface
}

func New() Config {
	return Config{}
}

func (c *Config) AddTarget(t TargetInterface) {
	c.Targets = append(c.Targets, t)
}

func (c *Config) AddTargets(t ...TargetInterface) {
	c.Targets = append(c.Targets, t...)
}
