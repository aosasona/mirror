package meta

import (
	"regexp"
)

type Optional int8

const (
	OptionalNone Optional = iota
	OptionalTrue
	OptionalFalse
)

func (o Optional) String() string {
	switch o {
	case OptionalTrue:
		return "true"
	case OptionalFalse:
		return "false"
	default:
		return "none"
	}
}
func (o Optional) IsOptional() bool { return o == OptionalTrue }
func (o Optional) True() bool       { return o == OptionalTrue }
func (o Optional) False() bool      { return o == OptionalFalse }
func (o Optional) None() bool       { return o == OptionalNone }

var FieldNameRegex = regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*$`)

type Meta struct {
	// OriginalName is the original name of the field in the Go struct
	OriginalName string

	// Name is the name of the field in the target language usually overridden by the user via parser hooks or struct tags
	Name string

	// Type is the type of the field in the target language usually overridden by the user via parser hooks or struct tags
	Type string

	// Optional is a flag indicating if the field is optional, depending on the target language, this may or may not be the same as nullable
	Optional Optional

	// Skip is a flag indicating if the field should be skipped during generation
	Skip bool
}
