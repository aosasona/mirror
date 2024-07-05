package jsonmeta

import (
	"reflect"
	"testing"

	"go.trulyao.dev/mirror/extractor/meta"
)

type TestStruct struct {
	Name         string `json:"first_name"`
	OmitEmpty    string `json:",omitempty"`
	ShouldSkip   string `json:"-"`
	Invalid      string `json:"invalid,omitempty,invalid"`
	Default      string `json:""`
	Tagless      string
	Formed       string `json:"formed,omitempty"`
	AsStr        string `json:",string"`
	Dash         string `json:"-,"`
	WithNameOnly string `json:"name,"`
}

var testStruct = reflect.TypeOf(TestStruct{})

func TestJSONTagParser_Parse(t *testing.T) {
	ok := true

	nameField, _ := testStruct.FieldByName("Name")
	omitemptyField, _ := testStruct.FieldByName("OmitEmpty")
	skipField, _ := testStruct.FieldByName("ShouldSkip")
	invalidField, _ := testStruct.FieldByName("Invalid")
	defaultField, _ := testStruct.FieldByName("Default")
	taglessField, _ := testStruct.FieldByName("Tagless")
	propertlyFormedField, _ := testStruct.FieldByName("Formed")
	asStrField, _ := testStruct.FieldByName("AsStr")
	dashField, _ := testStruct.FieldByName("Dash")
	withNameOnlyField, _ := testStruct.FieldByName("WithNameOnly")

	if !ok {
		panic("field not found")
	}

	tests := []struct {
		Name     string
		Source   reflect.StructField
		Expected *meta.Meta
		WantErr  bool
	}{
		{
			Name:   "properly parse tag with custom name",
			Source: nameField,
			Expected: &meta.Meta{
				OriginalName: "Name",
				Name:         "first_name",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "should be optional",
			Source: omitemptyField,
			Expected: &meta.Meta{
				OriginalName: "OmitEmpty",
				Name:         "OmitEmpty",
				Skip:         false,
				Optional:     true,
			},
		},
		{
			Name:   "skip tag",
			Source: skipField,
			Expected: &meta.Meta{
				OriginalName: "ShouldSkip",
				Name:         "ShouldSkip",
				Skip:         true,
				Optional:     false,
			},
		},
		{
			Name:    "fail to parse invalid tag",
			Source:  invalidField,
			WantErr: true,
		},
		{
			Name:   "produce tag with default name",
			Source: defaultField,
			Expected: &meta.Meta{
				OriginalName: "Default",
				Name:         "Default",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly handle tagless fields",
			Source: taglessField,
			Expected: &meta.Meta{
				OriginalName: "Tagless",
				Name:         "Tagless",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly handle properly formed tag",
			Source: propertlyFormedField,
			Expected: &meta.Meta{
				OriginalName: "Formed",
				Name:         "formed",
				Skip:         false,
				Optional:     true,
			},
		},
		{
			Name:   "properly handle string tag",
			Source: asStrField,
			Expected: &meta.Meta{
				OriginalName: "AsStr",
				Name:         "AsStr",
				Skip:         false,
				Optional:     false,
				Type:         "string",
			},
		},
		{
			Name:   "properly handle dash as name tag",
			Source: dashField,
			Expected: &meta.Meta{
				OriginalName: "Dash",
				Name:         "-",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly handle tag with name only and empty second pair",
			Source: withNameOnlyField,
			Expected: &meta.Meta{
				OriginalName: "WithNameOnly",
				Name:         "name",
				Skip:         false,
				Optional:     false,
			},
		},
	}

	for _, test := range tests {
		got, err := Extract(test.Source, nil)
		if err != nil && !test.WantErr {
			t.Errorf("[%s] wanted NO error, got error: %v", test.Name, err)
		}

		if err == nil && test.WantErr {
			t.Errorf("[%s] wanted error, got no error", test.Name)
		}

		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf("`%s`: expected %+v, got %+v", test.Name, test.Expected, got)
		}
	}
}
