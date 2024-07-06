package config

import (
	"fmt"

	"go.trulyao.dev/mirror/helper"
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
	Enabled               *bool    // can be used to disable or enable the generation of types, defaults to false
	Targets               []Target // the targets to generate types for
	UseTypeForObjects     *bool    // if true, will use `type Foo = ...` instead of `interface Foo {...}`, defaults to true
	ExpandObjectTypes     *bool    // if true, will expand object types instead of just using the name (e.g foo: { bar: string } instead of foo: Bar)
	PreferUnknown         *bool    // if true, will prefer unknown over any
	AllowUnexportedFields *bool    // if true, will include private fields

	IndentationType Indentation // the type of indentation to use, defaults to `Space`
	SpaceCount      *int        // the number of spaces to use for indentation, defaults to 4

	// TODO: implement custom types
	CustomTypes map[string]string // custom types to be used in the generation of types
}

func DefaultConfig() Config {
	return Config{
		Enabled: helper.Bool(true),
		Targets: []Target{
			{Filename: "generated.ts", PathToDirectory: ".", Language: LangTypescript},
		},
		UseTypeForObjects:     helper.Bool(true),
		ExpandObjectTypes:     helper.Bool(true),
		PreferUnknown:         helper.Bool(false),
		AllowUnexportedFields: helper.Bool(false),

		IndentationType: Space,
		SpaceCount:      helper.Int(4),
	}
}

func (c Config) EnabledOrDefault() bool {
	if c.Enabled == nil {
		if Debug {
			fmt.Println("c.Enabled is nil, check your config")
		}
		return false
	}

	return *c.Enabled
}

func (c Config) ExpandObjectTypesOrDefault() bool {
	if c.ExpandObjectTypes == nil {
		return false
	}

	return *c.ExpandObjectTypes
}

func (c Config) UseTypeForObjectsOrDefault() bool {
	if c.UseTypeForObjects == nil {
		return true
	}

	return *c.UseTypeForObjects
}

func (c Config) PreferUnknownOrDefault() bool {
	if c.PreferUnknown == nil {
		return false
	}

	return *c.PreferUnknown
}

func (c Config) AllowUnexportedFieldsOrDefault() bool {
	if c.AllowUnexportedFields == nil {
		return false
	}

	return *c.AllowUnexportedFields
}

func (c Config) IndentationTypeOrDefault() Indentation {
	if c.IndentationType == emptyIndent {
		return Space
	}

	return c.IndentationType
}

func (c Config) SpaceCountOrDefault() int {
	if c.SpaceCount == nil {
		return 4
	}

	return *c.SpaceCount
}

func (c Config) GetIndentation() string {
	indentationType := c.IndentationTypeOrDefault()
	spaceCount := c.SpaceCountOrDefault()

	if indentationType == Space {
		space := ""
		for i := 0; i < spaceCount; i++ {
			space += " "
		}

		return space
	}

	tab := ""
	repeatTabCount := int(spaceCount / 4)

	for i := 0; i < repeatTabCount; i++ {
		tab += "\t"
	}

	return tab
}

func (c Config) TargetsOrDefault() []Target {
	if c.Targets == nil {
		return []Target{}
	}

	return c.Targets
}

func (c Config) Merge(other Config) Config {
	enabled := c.EnabledOrDefault()
	targets := c.TargetsOrDefault()

	useTypeForObjects := c.UseTypeForObjectsOrDefault()
	expandObjectTypes := c.ExpandObjectTypesOrDefault()
	allowUnexportedFields := c.AllowUnexportedFieldsOrDefault()
	preferUnknown := c.PreferUnknownOrDefault()

	indentationType := c.IndentationTypeOrDefault()
	spaceCount := c.SpaceCountOrDefault()

	if other.Enabled != nil {
		enabled = *other.Enabled
	}

	if other.Targets != nil {
		targets = append(targets, other.Targets...)

		// Filter out duplicates
		targetMap := make(map[string]Target)
		for _, target := range targets {
			targetMap[target.Filename] = target

			if Debug {
				fmt.Println("Added target", target.Filename)
			}
		}

		targets = make([]Target, 0, len(targetMap))

		for _, target := range targetMap {
			targets = append(targets, target)
		}
	}

	if other.UseTypeForObjects != nil {
		useTypeForObjects = *other.UseTypeForObjects
	}

	if other.ExpandObjectTypes != nil {
		expandObjectTypes = *other.ExpandObjectTypes
	}

	if other.AllowUnexportedFields != nil {
		allowUnexportedFields = *other.AllowUnexportedFields
	}

	if other.PreferUnknown != nil {
		preferUnknown = *other.PreferUnknown
	}

	if other.IndentationType != emptyIndent {
		indentationType = other.IndentationType
	}

	if other.SpaceCount != nil {
		spaceCount = *other.SpaceCount
	}

	return Config{
		Enabled: &enabled,
		Targets: targets,

		UseTypeForObjects:     &useTypeForObjects,
		ExpandObjectTypes:     &expandObjectTypes,
		AllowUnexportedFields: &allowUnexportedFields,
		PreferUnknown:         &preferUnknown,

		IndentationType: indentationType,
		SpaceCount:      &spaceCount,
	}
}
