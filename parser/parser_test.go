package parser

import (
	"reflect"
	"testing"
)

type Test struct {
	Description string
	Source      any
	Expected    Item
	WantErr     bool
}

func Test_ParseItem(t *testing.T) {
	type (
		Foo   int
		Foo8  int8
		Foo16 int16
		Foo32 int32
		Foo64 int64
	)

	tests := []Test{
		{
			Description: "parse integer",
			Source:      *new(Foo),
			Expected:    Scalar{"Foo", TypeInteger, false},
		},
		{
			Description: "parse i8",
			Source:      *new(Foo8),
			Expected:    Scalar{"Foo8", TypeInteger, false},
		},
		{
			Description: "parse i16",
			Source:      *new(Foo16),
			Expected:    Scalar{"Foo16", TypeInteger, false},
		},
		{
			Description: "parse i32",
			Source:      *new(Foo32),
			Expected:    Scalar{"Foo32", TypeInteger, false},
		},
		{
			Description: "parse i64",
			Source:      *new(Foo64),
			Expected:    Scalar{"Foo64", TypeInteger, false},
		},
	}

	for _, tt := range tests {
		got, err := ParseItem(reflect.TypeOf(tt.Source))
		if err != nil && !tt.WantErr {
			t.Errorf("wanted NO error, got error `%s`", tt.Description)
		}

		if err == nil && tt.WantErr {
			t.Errorf("wanted error, got no error in `%s`", tt.Description)
		}

		if !reflect.DeepEqual(got, tt.Expected) {
			t.Errorf("wanted %v, got %v", tt.Expected, got)
		}
	}
}
