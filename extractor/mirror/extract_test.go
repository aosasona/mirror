package mirrormeta_test

import (
	"reflect"
	"testing"
	"time"

	"go.trulyao.dev/mirror/v2/extractor/meta"
	mirrormeta "go.trulyao.dev/mirror/v2/extractor/mirror"
)

type TestStruct struct {
	Name        string
	LastName    string    `mirror:"name:last_name"`
	Invalid     string    `mirror:"name:,random_opt"`
	Phone       string    `mirror:"name:phone_number,optional:true"`
	BMI         string    `mirror:"-"`
	NextOfKin   string    `mirror:"name:next_of_kin,skip:true"`
	Connections []string  `mirror:"name:connected_ids, type:Array<string>, optional:true"`
	Meta        any       `mirror:"name:meta, type:{'foo': string},"`
	CreatedAt   time.Time `mirror:"type:Date,skip:true,optional:true"`
}

var testStruct = reflect.TypeOf(TestStruct{})

func TestJSONTagParser_Parse(t *testing.T) {
	nameField, _ := testStruct.FieldByName("Name")
	lastNameField, _ := testStruct.FieldByName("LastName")
	invalidField, _ := testStruct.FieldByName("Invalid")
	phoneField, _ := testStruct.FieldByName("Phone")
	bmiField, _ := testStruct.FieldByName("BMI")
	nextOfKinField, _ := testStruct.FieldByName("NextOfKin")
	connectionsField, _ := testStruct.FieldByName("Connections")
	metaField, _ := testStruct.FieldByName("Meta")
	createAtField, _ := testStruct.FieldByName("CreatedAt")

	tests := []struct {
		Name     string
		Source   reflect.StructField
		Expected *meta.Meta
		WantErr  bool
	}{
		{
			Name:   "parse with missing mirror tag",
			Source: nameField,
			Expected: &meta.Meta{
				OriginalName: "Name",
				Name:         "Name",
				Skip:         false,
				Optional:     meta.OptionalNone,
			},
		},
		{
			Name:   "properly parse tag using v2.0 mirror tag",
			Source: lastNameField,
			Expected: &meta.Meta{
				OriginalName: "LastName",
				Name:         "last_name",
				Skip:         false,
				Optional:     meta.OptionalNone,
			},
		},
		{
			Name:   "gracefully handle invalid tag (kv pair)",
			Source: invalidField,
			Expected: &meta.Meta{
				OriginalName: "Invalid",
				Name:         "Invalid",
				Skip:         false,
				Optional:     meta.OptionalNone,
			},
		},
		{
			Name:   "parse name and optional props",
			Source: phoneField,
			Expected: &meta.Meta{
				OriginalName: "Phone",
				Name:         "phone_number",
				Skip:         false,
				Optional:     meta.OptionalTrue,
			},
		},
		{
			Name:   "skip tag using -",
			Source: bmiField,
			Expected: &meta.Meta{
				OriginalName: "BMI",
				Name:         "BMI",
				Skip:         true,
				Optional:     meta.OptionalNone,
			},
		},
		{
			Name:   "skip tag using skip:true",
			Source: nextOfKinField,
			Expected: &meta.Meta{
				OriginalName: "NextOfKin",
				Name:         "next_of_kin",
				Skip:         true,
				Optional:     meta.OptionalNone,
			},
		},
		{
			Name:   "parse all props and override type (with whitespace) - optional",
			Source: connectionsField,
			Expected: &meta.Meta{
				OriginalName: "Connections",
				Name:         "connected_ids",
				Skip:         false,
				Optional:     meta.OptionalTrue,
				Type:         "Array<string>",
			},
		},
		{
			Name:   "parse object type override",
			Source: metaField,
			Expected: &meta.Meta{
				OriginalName: "Meta",
				Name:         "meta",
				Skip:         false,
				Optional:     meta.OptionalNone,
				Type:         "{'foo': string}",
			},
		},
		{
			Name:   "parse type and skip",
			Source: createAtField,
			Expected: &meta.Meta{
				OriginalName: "CreatedAt",
				Name:         "CreatedAt",
				Skip:         true,
				Optional:     meta.OptionalTrue,
				Type:         "Date",
			},
		},
	}

	for _, test := range tests {
		got, err := mirrormeta.Extract(test.Source, nil)

		if err != nil && !test.WantErr {
			t.Errorf("[%s] wanted NO error, got error: %v", test.Name, err)
		}

		if err == nil && test.WantErr {
			t.Errorf("[%s] wanted error, got no error", test.Name)
		}

		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf(
				"failed to run case `%s`: expected %+v, got %+v",
				test.Name,
				test.Expected,
				got,
			)
		}
	}
}
