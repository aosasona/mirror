package mirrormeta

import (
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/v2/extractor/meta"
	"go.trulyao.dev/mirror/v2/helper"
)

func Extract(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	var (
		fieldMeta *meta.Meta
		opts      = make(map[string]string)
	)

	if root == nil {
		fieldMeta = new(meta.Meta)
	} else {
		fieldMeta = root
	}

	fieldMeta.OriginalName = helper.WithDefaultString(fieldMeta.OriginalName, field.Name)
	fieldMeta.Name = helper.WithDefaultString(fieldMeta.Name, field.Name)

	if fieldMeta.Name == "" {
		fieldMeta.Name = fieldMeta.OriginalName
	}

	mirrorTag := strings.TrimSpace(field.Tag.Get("mirror"))

	if mirrorTag == "" {
		return fieldMeta, nil
	}

	if mirrorTag == "-" {
		fieldMeta.Skip = true
		return fieldMeta, nil
	}

	tagFields := strings.Split(mirrorTag, ",")
	if len(tagFields) == 0 {
		return fieldMeta, nil
	}

	// split the props into key-value pairs
	for _, f := range tagFields {
		kv := strings.SplitN(f, ":", 2)
		if len(kv) != 2 {
			continue
		}
		opts[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	if skip, ok := opts["skip"]; ok {
		if helper.StringToBool(skip) {
			fieldMeta.Skip = true

			// if the tag is set to skip, then we don't need to parse the rest of the tag
			return fieldMeta, nil
		}
	}

	if name, ok := opts["name"]; ok {
		fieldMeta.Name = helper.WithDefaultString(strings.TrimSpace(name), field.Name)
	}

	if optional, ok := opts["optional"]; ok {
		// check if the optional tag is set to true or 1 (for backwards compatibility)
		if helper.StringToBool(optional) {
			fieldMeta.Optional = true
		}
	}

	if overrideType, ok := opts["type"]; ok {
		fieldMeta.Type = strings.TrimSpace(overrideType)
	}

	return fieldMeta, nil
}
