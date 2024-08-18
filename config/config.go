package config

import (
	"log/slog"
	"slices"

	"go.trulyao.dev/mirror/v2/types"
)

type Indentation int

const (
	indentNone Indentation = iota
	IndentSpace
	IndentTab
)

// Pointers have been used here to make sure the user actually sets the value and not just uses the default value
type Config struct {
	// Enabled can be used to disable or enable the generation of types, defaults to false
	Enabled bool

	// EnableParserCache will enable caching of parsed types
	EnableParserCache bool

	// Targets are the languages and files to generate types for, at least ONE target MUST be defined
	Targets []types.TargetInterface

	// FlattenEmbeddedStructs will flatten embedded structs into the parent struct
	//
	// For example:
	//
	// type Bar struct {
	//     BarField string
	// }
	//
	// type Foo struct {
	//     Bar
	// }
	//
	// will become:
	//
	// type Foo struct {
	//     BarField string
	// }
	//
	FlattenEmbeddedStructs bool
}

// DefaultConfig returns a new Config with default values (Mirror is disabled by default)
func DefaultConfig() *Config {
	return &Config{
		Enabled:                false,
		EnableParserCache:      true,
		FlattenEmbeddedStructs: true,
	}
}

// AddTarget() adds a target to the list of targets
func (c *Config) AddTarget(target types.TargetInterface) *Config {
	// Targets cannot be empty
	if target == nil {
		return c
	}

	// Check if the target is already in the list
	if slices.ContainsFunc(c.Targets, target.IsEquivalent) {
		slog.Warn(
			"target already exists in the list",
			slog.String("target", target.Name()), slog.String("path", target.Path()),
		)
		return c
	}

	c.Targets = append(c.Targets, target)
	return nil
}

// AddTargets() adds multiple targets to the list of targets
func (c *Config) AddTargets(t ...types.TargetInterface) {
	c.Targets = append(c.Targets, t...)
}
