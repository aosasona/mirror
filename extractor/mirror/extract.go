package mirrormeta

import (
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/v2/extractor/meta"
	"go.trulyao.dev/mirror/v2/helper"
)

func Extract(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	var fieldMeta *meta.Meta

	if root == nil {
		fieldMeta = new(meta.Meta)
	} else {
		fieldMeta = root
	}

	fieldMeta.OriginalName = helper.WithDefaultString(fieldMeta.OriginalName, field.Name)
	fieldMeta.Name = helper.WithDefaultString(fieldMeta.Name, field.Name)

	mirrorTag := strings.TrimSpace(field.Tag.Get("mirror"))

	if mirrorTag == "" {
		return fieldMeta, nil
	}

	if mirrorTag == "-" {
		fieldMeta.Skip = true
		return fieldMeta, nil
	}

	// Parse the mirror tag
	// NOTE: This required because of the very nature of how types are, they can contain almost anything and will be near impossible to simply split on or extract with regex. The parser also gives us more control in the long-term.
	p := NewMetaParser(mirrorTag)
	parsedMeta, err := p.Parse()
	if err != nil {
		return nil, err
	}

	// If the skip flag is set, then we don't need to parse the rest of the tag
	if parsedMeta.Skip != nil {
		fieldMeta.Skip = *parsedMeta.Skip
	}

	// Update the field meta with the parsed meta
	if parsedMeta.Name != nil {
		fieldMeta.Name = *parsedMeta.Name
	}

	if parsedMeta.Optional != nil {
		fieldMeta.Optional = *parsedMeta.Optional
	}

	if parsedMeta.Type != nil {
		fieldMeta.Type = *parsedMeta.Type
	}

	return fieldMeta, nil
}
