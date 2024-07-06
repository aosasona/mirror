package parser

import (
	"fmt"
	"reflect"
	"regexp"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/extractor"
	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
)

var ScalarRegex = regexp.MustCompile("^(int(8|16|32|64)|float(32|64)|string|bool)$")

type Options struct {
	OverrideNullable *bool
}

type Parser struct {
	config *config.Config
}

func New(config *config.Config) *Parser {
	return &Parser{config: config}
}

// Converts a type to an `Item` type that can be passed to generators
// Non-scalar types like classes (structs) and slices are expanded to include their root type
//
// ## Example
//
// ```go
// package foo
//
// import (
//
//	"go.trulyao.dev/mirror"
//	"go.trulyao.dev/mirror/parser"
//
// )
//
//	func main() {
//		   type Bar struct {
//			      Baz string
//		   }
//
//		  parser.Parse(reflect.TypeOf(Bar{}), parser.Options{
//	       OverrideNullable: mirror.Bool(false)
//	   })
//	}
//
// ```
func (p *Parser) Parse(source reflect.Type, opts ...Options) (Item, error) {
	opt := Options{}

	if len(opts) > 0 {
		if len(opts) > 1 {
			return nil, fmt.Errorf(
				"expected only one instance of `ParseOptions` passed to this function, got %d",
				len(opts),
			)
		}

		opt = opts[0]
	}

	nullable := false
	if opt.OverrideNullable != nil {
		nullable = *opt.OverrideNullable
	}

	switch source.Kind() {

	case reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Int,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint:

		return Scalar{source.Name(), TypeInteger, nullable}, nil

	case reflect.Float32, reflect.Float64:
		return Scalar{source.Name(), TypeFloat, nullable}, nil

	case reflect.String:
		return Scalar{source.Name(), TypeString, nullable}, nil

	case reflect.Bool:
		return Scalar{source.Name(), TypeBoolean, nullable}, nil

	case reflect.Map:
		keyItem, err := p.Parse(source.Key())
		if err != nil {
			return Map{}, err
		}

		valueItem, err := p.Parse(source.Elem())
		if err != nil {
			return Map{}, err
		}

		return Map{source.Name(), keyItem, valueItem, nullable}, nil

	case reflect.Struct:
		fields := make([]Field, 0)

		for i := 0; i < source.NumField(); i++ {
			field := source.Field(i)

			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			meta, err := p.ParseField(field)
			if err != nil {
				return nil, err
			}

			item, err := p.Parse(field.Type)
			if err != nil {
				return nil, err
			}

			fields = append(fields, Field{Name: meta.Name, BaseItem: item, Meta: meta})
		}

		return Struct{Name: source.Name(), Fields: fields, Nullable: nullable}, nil

	case reflect.Pointer:
		baseType := source.Elem().Kind()

		// We need to check it is a pointer a scalar type, so that we can abort and return a normal item marked as nullable
		if matched := ScalarRegex.Match([]byte(baseType.String())); matched {
			return p.Parse(source.Elem(), Options{
				OverrideNullable: helper.Bool(true),
			})
		}

		// Whatever we have left after filtering out scalar values is one of: array, slice or struct

		// If it is a struct, we need to expand it
		if source.Elem().Kind() == reflect.Struct {
			return p.Parse(source.Elem(), Options{
				OverrideNullable: helper.Bool(true),
			})
		}

		// TODO: handle other pointer types (array, slices, uintptr)
	}

	return nil, fmt.Errorf("not implemented for %s", source.Name())
}

func (p *Parser) ParseField(field reflect.StructField) (meta.Meta, error) {
	rootMeta := meta.Meta{}

	rootMeta.OriginalName = helper.WithDefaultString(rootMeta.OriginalName, field.Name)
	rootMeta.Name = helper.WithDefaultString(rootMeta.Name, field.Name)

	// Parse the JSON struct tag first
	jsonMeta, err := extractor.ExtractJSONMeta(field, &rootMeta)
	if err != nil {
		return meta.Meta{}, err
	}

	// Parse the custom `mirror` and `ts` struct tags to override the JSON struct tag if present
	mirrorMeta, err := extractor.ExtractMirrorMeta(field, jsonMeta)
	if err != nil {
		return meta.Meta{}, err
	}

	return *mirrorMeta, nil
}
