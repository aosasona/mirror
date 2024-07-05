package jsonmeta

import (
	"fmt"
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
)

func Extract(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	var fieldMeta *meta.Meta

	if root != nil {
		fieldMeta = root
	} else {
		fieldMeta = new(meta.Meta)
	}

	fieldMeta.OriginalName = helper.WithDefaultString(fieldMeta.OriginalName, field.Name)
	fieldMeta.Name = helper.WithDefaultString(fieldMeta.Name, field.Name)

	if fieldMeta.Name == "" {
		fieldMeta.Name = fieldMeta.OriginalName
	}

	jsonTag := strings.TrimSpace(field.Tag.Get("json"))

	// If the JSON tag is empty, return the fieldMeta as is so that we can proceed to the next tag if any
	if jsonTag == "" {
		return fieldMeta, nil
	}

	// Check for the presence of a dash (-) in the JSON tag
	if jsonTag == "-" {
		fieldMeta.Skip = true
		return fieldMeta, nil
	}

	// Parse with respect to the DEFAULT order, usually, the first field is the name and the second is to indicate if it should be optional or some other directive like string
	values := strings.Split(jsonTag, ",")
	length := len(values)

	if length == 0 {
		return fieldMeta, nil
	}

	// Check if the JSON tag is properly formatted, maximum of 2 values
	if length >= 3 {
		return nil, fmt.Errorf("invalid JSON tag: %s, too many values", jsonTag)
	}

	name := strings.TrimSpace(values[0])
	secondDirective := ""
	if length == 2 {
		secondDirective = strings.TrimSpace(values[1])
	}

	// Validate the JSON tag name, making sure it is not a `,omitempty` as that is valid but doesn't count as the name
	// Also `-` is a valid name if it is written as `json:"-,"`
	if name != "" {
		if !meta.FieldNameRegex.MatchString(name) && name != "-" {
			return nil, fmt.Errorf("invalid JSON tag name: %s", name)
		}
		fieldMeta.Name = name
	}

	// Validate the second directive
	switch secondDirective {
	case "omitempty":
		fieldMeta.Optional = true
	case "string":
		fieldMeta.Type = "string"
	case "":
		return fieldMeta, nil
	default:
		return nil, fmt.Errorf(
			"invalid JSON tag directive: %s, second item is not supported",
			secondDirective,
		)
	}

	return fieldMeta, nil
}
