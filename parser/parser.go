package parser

import (
	"reflect"

	"go.trulyao.dev/mirror/extractor"
	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
)

// Converts a type to an `Item` type that can be passed to generators
// Non-scalar types like classes (structs) and slices are expanded to include their root type
func ParseItem(source reflect.Type) (Item, error) {
	switch source.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return Scalar{source.Name(), TypeInteger, false}, nil
	case reflect.Float32, reflect.Float64:
		return Scalar{source.Name(), TypeFloat, false}, nil
	case reflect.String:
		return Scalar{source.Name(), TypeString, false}, nil
	case reflect.Bool:
		return Scalar{source.Name(), TypeBoolean, false}, nil
	case reflect.Map:
		keyItem, err := ParseItem(source.Key())
		if err != nil {
			return Map{}, err
		}

		valueItem, err := ParseItem(source.Elem())
		if err != nil {
			return Map{}, err
		}

		return Map{source.Name(), keyItem, valueItem}, nil
	}

	return nil, nil
}

func ParseField(field reflect.StructField) (*meta.Meta, error) {
	tag := new(meta.Meta)

	tag.OriginalName = helper.WithDefaultString(tag.OriginalName, field.Name)
	tag.Name = helper.WithDefaultString(tag.Name, field.Name)

	// Parse the JSON struct tag first
	if _, err := extractor.ExtractJSONMeta(field, tag); err != nil {
		return nil, err
	}

	// Parse the custom `mirror` and `ts` struct tags to override the JSON struct tag if present
	if _, err := extractor.ExtractMirrorMeta(field, tag); err != nil {
		return nil, err
	}

	return tag, nil
}
