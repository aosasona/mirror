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

type ParserInterface interface {
	AddSource(reflect.Type) error
	AddSources(...reflect.Type) error
	ParseN(int) (Item, error)

	Next() (Item, error)
	Done() bool
}

type Parser struct {
	config  *config.Config
	sources []reflect.Type
}

func New(config *config.Config) *Parser {
	return &Parser{config: config}
}

func (p *Parser) Done() bool {
	return len(p.sources) == 0
}

func (p *Parser) AddSource(source reflect.Type) error {
	if source == nil {
		return fmt.Errorf("source cannot be nil")
	}

	p.sources = append(p.sources, source)
	return nil
}

func (p *Parser) AddSources(sources ...reflect.Type) error {
	for _, source := range sources {
		if err := p.AddSource(source); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) Next() (Item, error) {
	if len(p.sources) == 0 {
		return nil, fmt.Errorf("no sources to parse")
	}

	source := p.sources[0]
	p.sources = p.sources[1:]

	return p.Parse(source)
}

// Parse the nth source in the list of sources (0-indexed)
func (p *Parser) ParseN(n int) (Item, error) {
	if n < 0 {
		return nil, fmt.Errorf("n must be a positive integer")
	}

	if len(p.sources) < n {
		return nil, fmt.Errorf("not enough sources to parse")
	}

	source := p.sources[n]

	return p.Parse(source)
}

// Converts a type to an `Item` type that can be passed to generators
// Non-scalar types like classes (structs) and slices are expanded to include their root type
//
// # This is used internally to convert a `reflect.Type` to a `parser.Item` type but is left exposed for advanced use cases where the user wants to convert a type to an `Item` type
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
		return p.parseMap(source, nullable)

	case reflect.Struct:
		return p.parseStruct(source, nullable)

	case reflect.Array, reflect.Slice:
		return p.parseList(source, nullable)

	case reflect.Func:
		return p.parseFunc(source, nullable)

	case reflect.Pointer:
		return p.Parse(source.Elem(), Options{
			OverrideNullable: helper.Bool(true),
		})

	case reflect.Interface:
		return p.parseInterface(source, nullable)
	}

	return nil, fmt.Errorf("not implemented for %s", source.Name())
}

func (p *Parser) parseField(field reflect.StructField) (meta.Meta, error) {
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

func (p *Parser) parseStruct(source reflect.Type, nullable bool) (Struct, error) {
	fields := make([]Field, 0)

	for i := 0; i < source.NumField(); i++ {
		field := source.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		meta, err := p.parseField(field)
		if err != nil {
			return Struct{}, err
		}

		item, err := p.Parse(field.Type)
		if err != nil {
			return Struct{}, err
		}

		fields = append(fields, Field{Name: meta.Name, BaseItem: item, Meta: meta})
	}

	return Struct{Name: source.Name(), Fields: fields, Nullable: nullable}, nil
}

func (p *Parser) parseMap(source reflect.Type, nullable bool) (Map, error) {
	keyItem, err := p.Parse(source.Key())
	if err != nil {
		return Map{}, err
	}

	valueItem, err := p.Parse(source.Elem())
	if err != nil {
		return Map{}, err
	}

	return Map{source.Name(), keyItem, valueItem, nullable}, nil
}

func (p *Parser) parseList(source reflect.Type, nullable bool) (List, error) {
	item, err := p.Parse(source.Elem())
	if err != nil {
		return List{}, err
	}

	length := EmptyLength
	if source.Kind() == reflect.Array {
		length = source.Len()
	}

	return List{Name: source.Name(), BaseItem: item, Nullable: nullable, Length: length}, nil
}

func (p *Parser) parseFunc(source reflect.Type, nullable bool) (Function, error) {
	params := make([]Item, 0)
	returns := make([]Item, 0)

	for i := 0; i < source.NumIn(); i++ {
		param, err := p.Parse(source.In(i))
		if err != nil {
			return Function{}, err
		}

		params = append(params, param)
	}

	for i := 0; i < source.NumOut(); i++ {
		ret, err := p.Parse(source.Out(i))
		if err != nil {
			return Function{}, err
		}

		returns = append(returns, ret)
	}

	return Function{
		Name:     source.Name(),
		Params:   params,
		Returns:  returns,
		Nullable: nullable,
	}, nil
}

func (p *Parser) parseInterface(source reflect.Type, nullable bool) (Item, error) {
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

var _ ParserInterface = &Parser{}
