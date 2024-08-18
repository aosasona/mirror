package parser

import (
	"reflect"

	"go.trulyao.dev/mirror/helper"
	jsonparser "go.trulyao.dev/mirror/parser/json"
	mirrorparser "go.trulyao.dev/mirror/parser/mirror"
	"go.trulyao.dev/mirror/parser/tag"
)

func Parse(field reflect.StructField) (*tag.Tag, error) {
	tag := new(tag.Tag)

	tag.OriginalName = helper.WithDefaultString(tag.OriginalName, field.Name)
	tag.Name = helper.WithDefaultString(tag.Name, field.Name)

	// Parse the JSON struct tag first
	if _, err := jsonparser.Parse(field, tag); err != nil {
		return nil, err
	}

	// Parse the custom `mirror` and `ts` struct tags to override the JSON struct tag if present
	if _, err := mirrorparser.Parse(field, tag); err != nil {
		return nil, err
	}

	return tag, nil
}
