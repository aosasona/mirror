package parser

import (
	"reflect"

	"go.trulyao.dev/mirror/extractor"
	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
)

func ParseItem(field reflect.Type) (Item, error) {
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
