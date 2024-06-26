package mirrorparser

import (
	"fmt"
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/helper"
	"go.trulyao.dev/mirror/parser/tag"
)

func Parse(field reflect.StructField, targetTag ...*tag.Tag) (*tag.Tag, error) {
	var (
		tag  = new(tag.Tag)
		opts = make(map[string]string)
	)

	if len(targetTag) > 0 {
		tag = targetTag[0]
	} else {
		tag.OriginalName = helper.WithDefaultString(tag.OriginalName, field.Name)
		tag.Name = helper.WithDefaultString(tag.Name, field.Name)
	}

	mirrorTag := strings.TrimSpace(field.Tag.Get("mirror"))

	// check if there is a `ts` tag in the field (for backwards compatibility)
	if mirrorTag == "" {
		mirrorTag = strings.TrimSpace(field.Tag.Get("ts"))
		if mirrorTag != "" {
			fmt.Println("[WARN:MIRROR] The legacy `ts` struct tag has been deprecated and will be removed in a future release")
		}
	}

	if mirrorTag == "" {
		return tag, nil
	}

	if mirrorTag == "-" {
		tag.Skip = true
		return tag, nil
	}

	tagFields := strings.Split(mirrorTag, ",")
	if len(tagFields) == 0 {
		return tag, nil
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
		if strings.TrimSpace(skip) == "true" || strings.TrimSpace(skip) == "1" {
			tag.Skip = true

			// if the tag is set to skip, then we don't need to parse the rest of the tag
			return tag, nil
		}
	}

	if name, ok := opts["name"]; ok {
		tag.Name = helper.WithDefaultString(strings.TrimSpace(name), field.Name)
	}

	if optional, ok := opts["optional"]; ok {
		// check if the optional tag is set to true or 1 (for backwards compatibility)
		if strings.TrimSpace(optional) == "true" || strings.TrimSpace(optional) == "1" {
			tag.Optional = true
		}
	}

	if overrideType, ok := opts["type"]; ok {
		tag.Type = strings.TrimSpace(overrideType)
	}

	return tag, nil
}
