package meta

import "regexp"

var FieldNameRegex = regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*$`)

type Meta struct {
	OriginalName string
	Name         string
	Type         string
	Optional     bool
	Skip         bool
}
