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
	TypeStruct  ItemType = "class"
	TypeList    ItemType = "list"
	TypeArray   ItemType = "array"
	TypeMap     ItemType = "map"
)

type Field struct {
	Name    string
	RawType ItemType
	Meta    *meta.Meta
}

type Item interface {
	ItemName() string
	ItemType() ItemType
}

// Represents a struct type
type Struct struct {
	Name   string
	Fields []Field
}

// Represents a scalar type like string, number, boolean, etc.
// But it also includes arrays and slices
type Scalar struct {
	Name string
	Type ItemType
}

// Represents a map type
type Pair struct {
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
func (p Pair) ItemName() string {
	return p.Name
}

func (p Pair) ItemType() ItemType {
	return TypeMap
}

var (
	_ Item = Scalar{}
	_ Item = Struct{}
	_ Item = Pair{}
)
