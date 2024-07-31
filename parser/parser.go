package parser

import (
	"fmt"
	"reflect"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/extractor"
	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
)

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

	case reflect.Array, reflect.Slice:
		item, err := p.Parse(source.Elem())
		if err != nil {
			return nil, err
		}

		length := EmptyLength
		if source.Kind() == reflect.Array {
			length = source.Len()
		}

		return List{Name: source.Name(), BaseItem: item, Nullable: nullable, Length: length}, nil

	case reflect.Func:
		params := make([]Item, 0)
		returns := make([]Item, 0)

		for i := 0; i < source.NumIn(); i++ {
			param, err := p.Parse(source.In(i))
			if err != nil {
				return nil, err
			}

			params = append(params, param)
		}

		for i := 0; i < source.NumOut(); i++ {
			ret, err := p.Parse(source.Out(i))
			if err != nil {
				return nil, err
			}

			returns = append(returns, ret)
		}

		return Function{
			Name:     source.Name(),
			Params:   params,
			Returns:  returns,
			Nullable: nullable,
		}, nil

	case reflect.Pointer:
		return p.Parse(source.Elem(), Options{
			OverrideNullable: helper.Bool(true),
		})

	case reflect.Interface:
		switch source.Name() {
		case "interface{}", "any":
			return Scalar{source.Name(), TypeAny, nullable}, nil
		case "error":
			return Scalar{source.Name(), TypeString, nullable}, nil
		case "time.Time":
			return Scalar{source.Name(), TypeDateTime, nullable}, nil
		default:
			return nil, fmt.Errorf("not implemented for %s", source.Name())
		}
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
