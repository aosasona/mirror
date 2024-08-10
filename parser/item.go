package parser

import (
	"go.trulyao.dev/mirror/extractor/meta"
)

type ItemType string

const (
	// Scalar types
	TypeInteger ItemType = "int"
	TypeFloat   ItemType = "float"
	TypeString  ItemType = "string"
	TypeBoolean ItemType = "bool"
	TypeAny     ItemType = "any"
	TypeByte    ItemType = "byte"

	TypeStruct ItemType = "struct"
	TypeList   ItemType = "list"
	TypeArray  ItemType = "array"
	TypeMap    ItemType = "map"

	TypeFunction  ItemType = "function"
	TypeTimestamp ItemType = "datetime"
)

// General interface to be adopted by anything that can or should be represented as an item
// NOTE: probably should be called node but I will come back later
type Item interface {
	ItemName() string
	ItemType() ItemType
}

// Representing a field in a struct
type Field struct {
	Name     string
	BaseItem Item
	Meta     meta.Meta
}

// Represents a struct type
type Struct struct {
	Name     string
	Fields   []Field
	Nullable bool
}

// Represents a scalar type like string, number, boolean, etc.
type Scalar struct {
	Name     string
	Type     ItemType
	Nullable bool
}

// Represents a list type; array or slice
const EmptyLength = -1 // used to represent a slice

type List struct {
	Name     string
	BaseItem Item
	Nullable bool
	Length   int // -1 if slice
}

// Represents a map type
type Map struct {
	Name     string
	Key      Item
	Value    Item
	Nullable bool
}

// Represents a function
type Function struct {
	Name     string
	Params   []Item
	Returns  []Item
	Nullable bool
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

// LIST
func (l List) ItemName() string {
	return l.Name
}

func (l List) ItemType() ItemType {
	return TypeList
}

func (l List) IsArray() bool {
	return l.Length != EmptyLength
}

// FUNCTION
func (f Function) ItemName() string {
	return f.Name
}

func (f Function) ItemType() ItemType {
	return TypeFunction
}

var (
	_ Item = Scalar{}
	_ Item = Struct{}
	_ Item = Map{}
	_ Item = List{}
	_ Item = Function{}
)
