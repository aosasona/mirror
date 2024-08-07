package config

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
	Targets []Target
}
