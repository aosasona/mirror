package parser

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	"go.trulyao.dev/mirror/v2/extractor"
	"go.trulyao.dev/mirror/v2/extractor/meta"
	"go.trulyao.dev/mirror/v2/helper"
)

type (
	Options struct {
		OverrideNullable bool
	}

	CacheValue struct {
		Options Options
		Item    *Item
	}

	// TODO: rename FlattenEmbeddedStructs to FlattenEmbeddedTypes
	Config struct {
		EnableCaching          bool
		FlattenEmbeddedStructs bool
	}

	OnParseItemFunc func(sourceName string, target Item) error

	OnParseFieldFunc func(parentType *reflect.Type, originalField *reflect.StructField, field *Field) error

	Parser struct {
		sources []reflect.Type
		cache   map[string]CacheValue

		enableCaching          bool
		flattenEmbeddedStructs bool

		// Hooks
		onParseItemFn  OnParseItemFunc
		onParseFieldFn OnParseFieldFunc
	}
)

// New creates a new parser
func New() *Parser {
	return &Parser{
		cache:                  make(map[string]CacheValue),
		sources:                []reflect.Type{},
		enableCaching:          true,
		flattenEmbeddedStructs: false,
	}
}

// Reset the parser to its initial state
func (p *Parser) Reset() {
	p.sources = make([]reflect.Type, 0)
	p.cache = make(map[string]CacheValue)
}

// Set the parser's configuration
func (p *Parser) SetConfig(config Config) error {
	p.enableCaching = config.EnableCaching
	p.flattenEmbeddedStructs = config.FlattenEmbeddedStructs

	return nil
}

// Lookup a source by name, returns the source and a boolean indicating if the source was found
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

// Get the sources to parse
func (p *Parser) Sources() []reflect.Type {
	return p.sources
}

// Set the sources to parse
func (p *Parser) SetSources(sources []reflect.Type) {
	p.sources = sources
}

// Enable or disable flattening of embedded structs
func (p *Parser) SetFlattenEmbeddedStructs(flatten bool) *Parser {
	p.flattenEmbeddedStructs = flatten
	return p
}

// Enable or disable caching
func (p *Parser) SetEnableCaching(enable bool) *Parser {
	p.enableCaching = enable
	return p
}

// Add a source to the parser
func (p *Parser) AddSource(source reflect.Type) error {
	if source == nil {
		return fmt.Errorf("source cannot be nil")
	}

	p.sources = append(p.sources, source)
	return nil
}

// Add multiple sources to the parser
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

// Sets the hook to run after parsing has been done
// This is run after the item has been parsed and is ready to be used (also pre-caching)
// If an error is returned, the item will not be cached or returned, and the error will be returned
func (p *Parser) OnParseItem(fn OnParseItemFunc) {
	if fn != nil {
		p.onParseItemFn = fn
	}
}

// Sets the hook to run after a field has been parsed
// This is run after the field has been parsed and is ready to be attached to the original item
// If an error is returned, the field will not be attached to the original item and the error will be returned
func (p *Parser) OnParseField(fn OnParseFieldFunc) {
	if fn != nil {
		p.onParseFieldFn = fn
	}
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
//	"go.trulyao.dev/mirror/v2"
//	"go.trulyao.dev/mirror/v2/parser"
//
// )
//
//	func main() {
//		   type Bar struct {
//			      Baz string
//		   }
//
//		  parser.Parse(reflect.TypeOf(Bar{}), parser.Options{
//	       OverrideNullable: false
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
		item = &Scalar{source.Name(), TypeInteger, nullable}

	case reflect.Float32, reflect.Float64:
		item = &Scalar{source.Name(), TypeFloat, nullable}

	case reflect.String:
		item = &Scalar{source.Name(), TypeString, nullable}

	case reflect.Bool:
		item = &Scalar{source.Name(), TypeBoolean, nullable}

	case reflect.Map:
		item, err = p.parseMap(source, nullable)

	case reflect.Struct:
		// Attempt to parse exempted structs like `sql.NullX` types
		item, err = p.parseExemptedStructs(source, nullable)
		if err != nil {
			// If it is not an exempted struct, parse it as a regular struct
			item, err = p.parseStruct(source, nullable)
		}

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

	// Add item to cache if caching is enabled
	if p.enableCaching {
		p.cache[cacheKey] = CacheValue{Options: opt, Item: &item}
	}

	// Run the `OnParseItem` hook if present
	if p.onParseItemFn != nil {
		if err := p.onParseItemFn(source.Name(), item); err != nil {
			return nil, err
		}
	}

	return item, nil
}

// Parse a struct field and extract the meta information
func (p *Parser) parseField(fieldName string, field reflect.StructField) (meta.Meta, error) {
	rootMeta := meta.Meta{}

	rootMeta.OriginalName = helper.WithDefaultString(fieldName, field.Name)
	rootMeta.Name = helper.WithDefaultString(fieldName, field.Name)

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

// Parse exempted structs like `sql.NullX` types and other built-in types
// TODO: add support for more exempted structs
func (p *Parser) parseExemptedStructs(source reflect.Type, nullable bool) (Item, error) {
	switch {
	case source == reflect.TypeOf(time.Time{}):
		return &Scalar{source.Name(), TypeTimestamp, nullable}, nil

	case source == reflect.TypeOf(time.Duration(0)):
		return &Scalar{source.Name(), TypeInteger, nullable}, nil

	case source == reflect.TypeOf([]byte{}):
		return &Scalar{source.Name(), TypeString, nullable}, nil

	case source == reflect.TypeOf([]interface{}{}):
		return &List{
			ItemName: source.Name(),
			BaseItem: &Scalar{"any", TypeAny, nullable},
			Nullable: nullable,
		}, nil

	// SQL types
	case source == reflect.TypeOf(sql.NullBool{}):
		return &Scalar{source.Name(), TypeBoolean, true}, nil

	case source == reflect.TypeOf(sql.NullFloat64{}):
		return &Scalar{source.Name(), TypeFloat, true}, nil

	case
		source == reflect.TypeOf(sql.NullInt64{}),
		source == reflect.TypeOf(sql.NullInt32{}),
		source == reflect.TypeOf(sql.NullInt16{}):
		return &Scalar{source.Name(), TypeInteger, true}, nil

	case source == reflect.TypeOf(sql.NullString{}):
		return &Scalar{source.Name(), TypeString, true}, nil

	case source == reflect.TypeOf(sql.NullTime{}):
		return &Scalar{source.Name(), TypeTimestamp, true}, nil

	case source == reflect.TypeOf(sql.NullByte{}):
		return &Scalar{source.Name(), TypeByte, true}, nil

	default:
		return nil, fmt.Errorf("not implemented for %s", source.Name())
	}
}

// Parse a struct type
func (p *Parser) parseStruct(source reflect.Type, nullable bool) (*Struct, error) {
	fields := make([]Field, 0)

	withOnParseFieldHook := func(field *Field, sourceField *reflect.StructField) error {
		if p.onParseFieldFn != nil {
			if err := p.onParseFieldFn(&source, sourceField, field); err != nil {
				return fmt.Errorf("failed to run `OnParseField` hook: %s", err.Error())
			}
		}

		return nil
	}

	for i := 0; i < source.NumField(); i++ {
		sourceField := source.Field(i)

		// Skip unexported fields
		if !sourceField.IsExported() {
			continue
		}

		// If it is embedded, parse it as part of the original struct (flatten it)
		if p.flattenEmbeddedStructs && sourceField.Anonymous &&
			sourceField.Type.Kind() == reflect.Struct {
			item, err := p.Parse(sourceField.Type)
			if err != nil {
				return &Struct{}, err
			}

			switch nestedItem := item.(type) {
			case *Struct:
				fields = append(fields, nestedItem.Fields...)

			default:
				// Account for the case where the embedded type is not a struct
				name := helper.WithDefaultString(nestedItem.Name(), sourceField.Name)
				nestedItemMeta, err := p.parseField(name, sourceField)
				if err != nil {
					return &Struct{}, fmt.Errorf("failed to parse embedded field `%s`: %s", sourceField.Name, err.Error())
				}

				nestedField := Field{ItemName: nestedItemMeta.Name, BaseItem: nestedItem, Meta: nestedItemMeta}
				if err := withOnParseFieldHook(&nestedField, &sourceField); err != nil {
					return &Struct{}, err
				}

				fields = append(fields, nestedField)
			}

			continue
		}

		meta, err := p.parseField("", sourceField)
		if err != nil {
			return &Struct{}, err
		}

		item, err := p.Parse(sourceField.Type)
		if err != nil {
			return &Struct{}, err
		}

		field := Field{ItemName: meta.Name, BaseItem: item, Meta: meta}
		if err := withOnParseFieldHook(&field, &sourceField); err != nil {
			return &Struct{}, err
		}

		fields = append(fields, field)
	}

	return &Struct{ItemName: source.Name(), Fields: fields, Nullable: nullable}, nil
}

// Parse a map type
func (p *Parser) parseMap(source reflect.Type, nullable bool) (*Map, error) {
	keyItem, err := p.Parse(source.Key())
	if err != nil {
		return &Map{}, err
	}

	valueItem, err := p.Parse(source.Elem())
	if err != nil {
		return &Map{}, err
	}

	return &Map{source.Name(), keyItem, valueItem, nullable}, nil
}

// Parse a list type (slice or array)
func (p *Parser) parseList(source reflect.Type, nullable bool) (*List, error) {
	item, err := p.Parse(source.Elem())
	if err != nil {
		return &List{}, err
	}

	length := EmptyLength
	if source.Kind() == reflect.Array {
		length = source.Len()
	}

	return &List{ItemName: source.Name(), BaseItem: item, Nullable: nullable, Length: length}, nil
}

// Parse a function type
func (p *Parser) parseFunc(source reflect.Type, nullable bool) (*Function, error) {
	params := make([]Item, 0)
	returns := make([]Item, 0)

	for i := 0; i < source.NumIn(); i++ {
		param, err := p.Parse(source.In(i))
		if err != nil {
			return &Function{}, err
		}

		params = append(params, param)
	}

	for i := 0; i < source.NumOut(); i++ {
		ret, err := p.Parse(source.Out(i))
		if err != nil {
			return &Function{}, err
		}

		returns = append(returns, ret)
	}

	return &Function{
		ItemName: source.Name(),
		Params:   params,
		Returns:  returns,
		Nullable: nullable,
	}, nil
}

// Parse an interface type
// This accounts for various types like `interface{}`, `error`, `time.Time`, `sql.NullX` types
func (p *Parser) parseInterface(source reflect.Type, nullable bool) (Item, error) {
	switch source.Name() {
	case "interface{}", "any":
		return &Scalar{source.Name(), TypeAny, nullable}, nil
	case "error":
		return &Scalar{source.Name(), TypeString, nullable}, nil
	default:
		return nil, fmt.Errorf("not implemented for %s", source.Name())
	}
}
