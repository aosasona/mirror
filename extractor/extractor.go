package extractor

import (
	"reflect"

	jsonmeta "go.trulyao.dev/mirror/v2/extractor/json"
	"go.trulyao.dev/mirror/v2/extractor/meta"
	mirrormeta "go.trulyao.dev/mirror/v2/extractor/mirror"
)

// ExtractJSONMeta extracts meta information from a field with the `json` tag
func ExtractJSONMeta(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	return jsonmeta.Extract(field, root)
}

// ExtractMirrorMeta extracts meta information from a field with the `mirror` tag
func ExtractMirrorMeta(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	return mirrormeta.Extract(field, root)
}
