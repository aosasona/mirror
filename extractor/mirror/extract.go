package mirrormeta

import (
	"fmt"
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/extractor/meta"
	"go.trulyao.dev/mirror/helper"
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

	// check if there is a `ts` tag in the field (for backwards compatibility)
	if mirrorTag == "" {
		// Deprecated: The `ts` struct tag has been deprecated and will be removed in a future release
		mirrorTag = strings.TrimSpace(field.Tag.Get("ts"))
		if mirrorTag != "" {
			fmt.Println(
				"[WARN:MIRROR] The legacy `ts` struct tag has been deprecated and will be removed in a future release",
			)
		}
	}

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
		kv := strings.Split(f, ":")
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
