package extractor

import (
	"reflect"

	jsonmeta "go.trulyao.dev/mirror/extractor/json"
	"go.trulyao.dev/mirror/extractor/meta"
	mirrormeta "go.trulyao.dev/mirror/extractor/mirror"
)

func ExtractJSONMeta(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	return jsonmeta.Extract(field, root)
}

func ExtractMirrorMeta(field reflect.StructField, root *meta.Meta) (*meta.Meta, error) {
	return mirrormeta.Extract(field, root)
}
