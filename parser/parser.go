package parser

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"go.trulyao.dev/mirror/extractor"
	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
)

type Options struct {
	OverrideNullable bool
}

type CacheValue struct {
	Options Options
	Item    *Item
}

type Parser struct {
	sources []reflect.Type
	cache   map[string]CacheValue

	enableCaching          bool
	flattenEmbeddedStructs bool
}

func New() *Parser {
	return &Parser{
		cache:                  make(map[string]CacheValue),
		sources:                []reflect.Type{},
		enableCaching:          true,
		flattenEmbeddedStructs: false,
	}
}

func (p *Parser) LookupByName(name string) (Item, bool) {
	for _, source := range p.sources {
		if source.Name() == name {
			item, err := p.Parse(source)
			if err != nil {
				return nil, false
			}

			return item, true
		}
	}

	return nil, false
}

func (p *Parser) Reset() {
	p.sources = make([]reflect.Type, 0)
}

func (p *Parser) Sources() []reflect.Type {
	return p.sources
}

func (p *Parser) SetSources(sources []reflect.Type) {
	p.sources = sources
}

func (p *Parser) SetFlattenEmbeddedStructs(flatten bool) *Parser {
	p.flattenEmbeddedStructs = flatten
	return p
}

func (p *Parser) SetEnableCaching(enable bool) *Parser {
	p.enableCaching = enable
	return p
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

// Count the number of sources left to parse
func (p *Parser) Count() int {
	return len(p.sources)
}

// Check if there are any sources left to parse
func (p *Parser) Done() bool {
	return len(p.sources) == 0
}

// Parse the next source in the list of sources, this function consumes the source and removes it from the list
// Call `Done` to check if there are any sources left
func (p *Parser) Next() (Item, error) {
	if len(p.sources) == 0 {
		return nil, fmt.Errorf("no sources to parse")
	}

	source := p.sources[0]
	p.sources = p.sources[1:]

	return p.Parse(source)
}

// Iterate over all sources and call the function `f` on each source
// Unlike `Next`, this function does not consume the sources and can be called multiple times
func (p *Parser) Iterate(f func(Item) error) error {
	for _, source := range p.sources {
		item, err := p.Parse(source)
		if err != nil {
			return err
		}

		if err := f(item); err != nil {
			return err
		}
	}

	return nil
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
	cacheKey := base64.StdEncoding.EncodeToString(
		[]byte(source.PkgPath() + ":" + source.Name()),
	)

	if p.enableCaching {
		if value, ok := p.cache[cacheKey]; ok && reflect.DeepEqual(value.Options, opts) {
			return *value.Item, nil
		}
	}

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
	if opt.OverrideNullable != nullable {
		nullable = opt.OverrideNullable
	}

	var (
		item Item
		err  error
	)

	switch source.Kind() {

	case
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Int,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint:
		item = Scalar{source.Name(), TypeInteger, nullable}

	case reflect.Float32, reflect.Float64:
		item = Scalar{source.Name(), TypeFloat, nullable}

	case reflect.String:
		item = Scalar{source.Name(), TypeString, nullable}

	case reflect.Bool:
		item = Scalar{source.Name(), TypeBoolean, nullable}

	case reflect.Map:
		item, err = p.parseMap(source, nullable)

	case reflect.Struct:
		item, err = p.parseStruct(source, nullable)

	case reflect.Array, reflect.Slice:
		item, err = p.parseList(source, nullable)

	case reflect.Func:
		item, err = p.parseFunc(source, nullable)

	case reflect.Pointer:
		item, err = p.Parse(source.Elem(), Options{OverrideNullable: true})

	case reflect.Interface:
		item, err = p.parseInterface(source, nullable)

	default:
		return nil, fmt.Errorf("not implemented for %s", source.Name())
	}

	if err != nil {
		return nil, err
	}

	if p.enableCaching {
		p.cache[cacheKey] = CacheValue{Options: opt, Item: &item}
	}

	return item, nil
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
		if !field.IsExported() {
			continue
		}

		// If it is embedded, parse it as part of the original struct (flatten it)
		if p.flattenEmbeddedStructs && field.Anonymous && field.Type.Kind() == reflect.Struct {
			item, err := p.Parse(field.Type)
			if err != nil {
				return Struct{}, err
			}

			if embeddedFields, ok := item.(Struct); ok {
				fields = append(fields, embeddedFields.Fields...)
			}

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

		fields = append(fields, Field{ItemName: meta.Name, BaseItem: item, Meta: meta})
	}

	return Struct{ItemName: source.Name(), Fields: fields, Nullable: nullable}, nil
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

	return List{ItemName: source.Name(), BaseItem: item, Nullable: nullable, Length: length}, nil
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
		ItemName: source.Name(),
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
		return Scalar{source.Name(), TypeTimestamp, nullable}, nil

	// There isn't a good way to force the parser to parse a `sql.NullX` type as a scalar instead of structs (users can have the same name for their structs and the `sql` package structs)
	// But if they are ever passed as interfaces _somehow_, we can handle them here
	case "sql.NullString":
		return Scalar{source.Name(), TypeString, true}, nil
	case "sql.NullInt64", "sql.NullInt32", "sql.NullInt16":
		return Scalar{source.Name(), TypeInteger, true}, nil
	case "sql.NullFloat64":
		return Scalar{source.Name(), TypeFloat, true}, nil
	case "sql.NullBool":
		return Scalar{source.Name(), TypeBoolean, true}, nil
	case "sql.NullTime":
		return Scalar{source.Name(), TypeTimestamp, true}, nil
	case "sql.NullByte":
		return Scalar{source.Name(), TypeByte, true}, nil
	default:
		return nil, fmt.Errorf("not implemented for %s", source.Name())
	}
}
