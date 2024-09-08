package parser

import (
	"go.trulyao.dev/mirror/v2/extractor/meta"
)

type Type string

const (
	// Scalar types
	TypeInteger   Type = "int"
	TypeFloat     Type = "float"
	TypeString    Type = "string"
	TypeBoolean   Type = "bool"
	TypeAny       Type = "any"
	TypeByte      Type = "byte"
	TypeTimestamp Type = "datetime"

	// Collection types
	TypeStruct Type = "struct"
	TypeList   Type = "list"
	TypeArray  Type = "array"
	TypeMap    Type = "map"

	TypeFunction Type = "function"
)

// General interface to be adopted by anything that can or should be represented as an item
// NOTE: probably should be called node but I will come back later
type Item interface {
	Name() string
	Type() Type
	IsScalar() bool
	IsNullable() bool
}

// Represents a field in a struct
type Field struct {
	ItemName string
	BaseItem Item
	Meta     meta.Meta
}

// Represents a struct type
type Struct struct {
	ItemName string
	Fields   []Field
	Nullable bool
}

// Represents a scalar type like string, number, boolean, etc.
type Scalar struct {
	ItemName string
	ItemType Type
	Nullable bool
}

// Represents a list type; array or slice
const EmptyLength = -1 // used to represent a slice

// Represents a list (array or slice) type
type List struct {
	ItemName string
	BaseItem Item
	Nullable bool
	Length   int // -1 if slice
}

// Represents a map type
type Map struct {
	ItemName string
	Key      Item
	Value    Item
	Nullable bool
}

// Represents a function
type Function struct {
	ItemName string
	Params   []Item
	Returns  []Item
	Nullable bool
}

// SCALAR
func (s *Scalar) Name() string {
	return s.ItemName
}

func (s *Scalar) Type() Type {
	return s.ItemType
}

func (s *Scalar) IsScalar() bool {
	return true
}

func (s *Scalar) IsNullable() bool {
	return s.Nullable
}

// STRUCT
func (s *Struct) Name() string {
	return s.ItemName
}

func (s *Struct) Type() Type {
	return TypeStruct
}

func (s *Struct) IsScalar() bool {
	return false
}

func (s *Struct) IsNullable() bool {
	return s.Nullable
}

// PAIR
func (m *Map) Name() string {
	return m.ItemName
}

func (m *Map) Type() Type {
	return TypeMap
}

func (m *Map) IsScalar() bool {
	return false
}

func (m *Map) IsNullable() bool {
	return m.Nullable
}

// LIST
func (l *List) Name() string {
	return l.ItemName
}

func (l *List) Type() Type {
	return TypeList
}

func (l *List) IsScalar() bool {
	return false
}

func (l *List) IsArray() bool {
	return l.Length != EmptyLength
}

func (l *List) IsNullable() bool {
	return l.Nullable
}

// FUNCTION
func (f *Function) Name() string {
	return f.ItemName
}

func (f *Function) Type() Type {
	return TypeFunction
}

func (f *Function) IsScalar() bool {
	return false
}

func (f *Function) IsNullable() bool {
	return f.Nullable
}

var (
	_ Item = (*Scalar)(nil)
	_ Item = (*Struct)(nil)
	_ Item = (*Map)(nil)
	_ Item = (*List)(nil)
	_ Item = (*Function)(nil)
)
