package config

import (
	"slices"

	"go.trulyao.dev/mirror/types"
)

type Indentation int

const (
	emptyIndent Indentation = iota
	Space
	Tab
)

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

func (c *Config) AddTarget(target types.TargetInterface) *Config {
	// Targets cannot be empty
	if c.Targets == nil {
		return c
	}

	// Check if the target is already in the list
	if slices.ContainsFunc(c.Targets, target.IsEquivalent) {
		return c
	}

	c.Targets = append(c.Targets, target)
	return nil
}

func (c *Config) AddTargets(t ...types.TargetInterface) {
	c.Targets = append(c.Targets, t...)
}
