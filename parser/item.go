package parser

import (
	"go.trulyao.dev/mirror/extractor/meta"
)

type ItemType string

const (
	TypeInteger ItemType = "int"
	TypeFloat   ItemType = "float"
	TypeString  ItemType = "string"
	TypeBoolean ItemType = "bool"
	TypeStruct  ItemType = "struct"
	TypeList    ItemType = "list"
	TypeArray   ItemType = "array"
	TypeMap     ItemType = "map"
)

// General interface to be adopted by anything that can or should be represented as an item
// NOTE: probably should be called node but I will come back later
type Item interface {
	ItemName() string
	ItemType() ItemType
}

// Representing a field in a struct
type Field struct {
	Name    string
	RawType ItemType
	Meta    *meta.Meta
}

// Represents a struct type
type Struct struct {
	Name   string
	Fields []Field
}

// Represents a scalar type like string, number, boolean, etc.
// But it also includes arrays and slices
type Scalar struct {
	Name     string
	Type     ItemType
	Nullable bool
}

// Represents a map type
type Map struct {
	Name  string
	Key   Item
	Value Item
}

// SCALAR
func (s Scalar) ItemName() string {
	return s.Name
}

func (s Scalar) ItemType() ItemType {
	return s.Type
}

// STRUCT
func (s Struct) ItemName() string {
	return s.Name
}

func (s Struct) ItemType() ItemType {
	return TypeStruct
}

// PAIR
func (p Map) ItemName() string {
	return p.Name
}

func (p Map) ItemType() ItemType {
	return TypeMap
}

var (
	_ Item = Scalar{}
	_ Item = Struct{}
	_ Item = Map{}
)
