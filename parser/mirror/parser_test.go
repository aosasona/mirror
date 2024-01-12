package mirrorparser

import (
	"reflect"
	"testing"

	"go.trulyao.dev/mirror/parser/tag"
)

type TestStruct struct {
	Name        string   `ts:"name:first_name"`
	LastName    string   `mirror:"name:last_name"`
	Invalid     string   `mirror:"name:,random_opt"`
	Phone       string   `mirror:"name:phone_number,optional:true"`
	BMI         string   `mirror:"-"`
	NextOfKin   string   `mirror:"name:next_of_kin,skip:true"`
	Connections []string `mirror:"name:connected_ids, type:Array<string>, optional:true"`
}

var testStruct = reflect.TypeOf(TestStruct{})

func TestJSONTagParser_Parse(t *testing.T) {
	ok := true

	nameField, ok := testStruct.FieldByName("Name")
	lastNameField, ok := testStruct.FieldByName("LastName")
	invalidField, ok := testStruct.FieldByName("Invalid")
	phoneField, ok := testStruct.FieldByName("Phone")
	bmiField, ok := testStruct.FieldByName("BMI")
	nextOfKinField, ok := testStruct.FieldByName("NextOfKin")
	connectionsField, ok := testStruct.FieldByName("Connections")

	if !ok {
		panic("field not found")
	}

	tests := []struct {
		Name     string
		Source   reflect.StructField
		Expected *tag.Tag
		WantErr  bool
	}{
		{
			Name:   "properly parse tag using v1.0 ts tag",
			Source: nameField,
			Expected: &tag.Tag{
				OriginalName: "Name",
				Name:         "first_name",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "properly parse tag using v2.0 mirror tag",
			Source: lastNameField,
			Expected: &tag.Tag{
				OriginalName: "LastName",
				Name:         "last_name",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "gracefully handle invalid tag (kv pair)",
			Source: invalidField,
			Expected: &tag.Tag{
				OriginalName: "Invalid",
				Name:         "Invalid",
				Skip:         false,
				Optional:     false,
			},
		},
		{
			Name:   "parse name and optional props",
			Source: phoneField,
			Expected: &tag.Tag{
				OriginalName: "Phone",
				Name:         "phone_number",
				Skip:         false,
				Optional:     true,
			},
		},
		{
			Name:   "skip tag using -",
			Source: bmiField,
			Expected: &tag.Tag{
				OriginalName: "BMI",
				Name:         "BMI",
				Skip:         true,
				Optional:     false,
			},
		},
		{
			Name:   "skip tag using skip:true",
			Source: nextOfKinField,
			Expected: &tag.Tag{
				OriginalName: "NextOfKin",
				Name:         "NextOfKin",
				Skip:         true,
				Optional:     false,
			},
		},
		{
			Name:   "parse all props and override type (with whitespace) - optional",
			Source: connectionsField,
			Expected: &tag.Tag{
				OriginalName: "Connections",
				Name:         "connected_ids",
				Skip:         false,
				Optional:     true,
				Type:         "Array<string>",
			},
		},
	}

	for _, test := range tests {
		got, err := Parse(test.Source)
		if (err != nil) != test.WantErr {
			t.Errorf("failed to run case `%s`: unexpected error: %v", test.Name, err)
		}

		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf("failed to run case `%s`: expected %+v, got %+v", test.Name, test.Expected, got)
		}
	}
}
