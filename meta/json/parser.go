package jsonmeta

import (
	"fmt"
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/helper"
	"go.trulyao.dev/mirror/meta"
)

func Parse(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	var fieldMeta *meta.Meta

	if root != nil {
		fieldMeta = root
	} else {
		fieldMeta = new(meta.Meta)
	}

	fieldMeta.OriginalName = helper.WithDefaultString(fieldMeta.OriginalName, field.Name)
	fieldMeta.Name = helper.WithDefaultString(fieldMeta.Name, field.Name)

	jsonTag := strings.TrimSpace(field.Tag.Get("json"))
	if jsonTag == "" {
		return fieldMeta, nil
	}

	if jsonTag == "-" {
		fieldMeta.Skip = true
		return fieldMeta, nil
	}

	// Parse with respect to the DEFAULT order, usually, the first field is the name and the second is to indicate if it should be optional
	values := strings.Split(jsonTag, ",")

	if len(values) == 0 {
		return fieldMeta, nil
	}

	name := strings.TrimSpace(values[0])
	if name != "" {
		// Validate the JSON tag name
		if !meta.FieldNameRegex.MatchString(name) {
			return nil, fmt.Errorf("invalid JSON tag name: %s", name)
		}

		fieldMeta.Name = name
	}

	if len(values) > 1 {
		if strings.TrimSpace(values[1]) == "omitempty" {
			fieldMeta.Optional = true
		} else if strings.TrimSpace(values[1]) == "string" {
			fieldMeta.Type = "string"
		} else {
			return nil, fmt.Errorf("invalid JSON tag value: %s", values[1])
		}
	}

	return fieldMeta, nil
}
