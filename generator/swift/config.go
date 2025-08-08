package swift

import "go.trulyao.dev/mirror/v2/types"

type Config struct{}

// AddCustomType implements types.TargetInterface.
func (c *Config) AddCustomType(string, string) {
	panic("unimplemented")
}

// Extension implements types.TargetInterface.
func (c *Config) Extension() string {
	panic("unimplemented")
}

// Generator implements types.TargetInterface.
func (c *Config) Generator() types.GeneratorInterface {
	panic("unimplemented")
}

// Header implements types.TargetInterface.
func (c *Config) Header() string {
	panic("unimplemented")
}

// ID implements types.TargetInterface.
func (c *Config) ID() string {
	panic("unimplemented")
}

// IsEquivalent implements types.TargetInterface.
func (c *Config) IsEquivalent(types.TargetInterface) bool {
	panic("unimplemented")
}

// Language implements types.TargetInterface.
func (c *Config) Language() string {
	panic("unimplemented")
}

// Name implements types.TargetInterface.
func (c *Config) Name() string {
	panic("unimplemented")
}

// Path implements types.TargetInterface.
func (c *Config) Path() string {
	panic("unimplemented")
}

// Prefix implements types.TargetInterface.
func (c *Config) Prefix() string {
	panic("unimplemented")
}

// Validate implements types.TargetInterface.
func (c *Config) Validate() error {
	panic("unimplemented")
}
