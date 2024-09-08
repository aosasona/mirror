package mirrormeta

import (
	"reflect"
	"testing"

	"go.trulyao.dev/mirror/v2/extractor/meta"
)

type TestStruct struct {
	Name        string   `ts:"name:first_name"`
	LastName    string   `                     mirror:"name:last_name"`
	Invalid     string   `                     mirror:"name:,random_opt"`
	Phone       string   `                     mirror:"name:phone_number,optional:true"`
	BMI         string   `                     mirror:"-"`
	NextOfKin   string   `                     mirror:"name:next_of_kin,skip:true"`
	Connections []string `                     mirror:"name:connected_ids, type:Array<string>, optional:true"`
}

var testStruct = reflect.TypeOf(TestStruct{})

func TestJSONTagParser_Parse(t *testing.T) {
	ok := true

	nameField, _ := testStruct.FieldByName("Name")
	lastNameField, _ := testStruct.FieldByName("LastName")
	invalidField, _ := testStruct.FieldByName("Invalid")
	phoneField, _ := testStruct.FieldByName("Phone")
	bmiField, _ := testStruct.FieldByName("BMI")
	nextOfKinField, _ := testStruct.FieldByName("NextOfKin")
	connectionsField, _ := testStruct.FieldByName("Connections")

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
			Name:   "properly parse tag using v1.0 ts tag",
			Source: nameField,
			Expected: &meta.Meta{
				OriginalName: "Name",
				Name:         "Name",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly parse tag using v2.0 mirror tag",
			Source: lastNameField,
			Expected: &meta.Meta{
				OriginalName: "LastName",
				Name:         "last_name",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "gracefully handle invalid tag (kv pair)",
			Source: invalidField,
			Expected: &meta.Meta{
				OriginalName: "Invalid",
				Name:         "Invalid",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "parse name and optional props",
			Source: phoneField,
			Expected: &meta.Meta{
				OriginalName: "Phone",
				Name:         "phone_number",
				Skip:         false,
				Optional:     true,
			},
		},
		{
			Name:   "skip tag using -",
			Source: bmiField,
			Expected: &meta.Meta{
				OriginalName: "BMI",
				Name:         "BMI",
				Skip:         true,
				Optional:     false,
			},
		},
		{
			Name:   "skip tag using skip:true",
			Source: nextOfKinField,
			Expected: &meta.Meta{
				OriginalName: "NextOfKin",
				Name:         "NextOfKin",
				Skip:         true,
				Optional:     false,
			},
		},
		{
			Name:   "parse all props and override type (with whitespace) - optional",
			Source: connectionsField,
			Expected: &meta.Meta{
				OriginalName: "Connections",
				Name:         "connected_ids",
				Skip:         false,
				Optional:     true,
				Type:         "Array<string>",
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
			t.Errorf(
				"failed to run case `%s`: expected %+v, got %+v",
				test.Name,
				test.Expected,
				got,
			)
		}
	}
}
