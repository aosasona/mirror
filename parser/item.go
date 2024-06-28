package parser

import "go.trulyao.dev/mirror/parser/tag"

type ItemType int

const (
	TypeInteger ItemType = iota
	TypeFloat
	TypeString
	TypeBoolean
	TypeClass
	TypeList
	TypeMap
)

type Field struct {
	Name    string
	RawType ItemType
	Tag     *tag.Tag
}

type Item interface {
	ItemName() string
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

// STRUCT
func (s Struct) ItemName() string {
	return s.Name
}

// PAIR
func (p Pair) ItemName() string {
	return p.Name
}

var (
	_ Item = Scalar{}
	_ Item = Struct{}
	_ Item = Pair{}
)
