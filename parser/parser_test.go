package parser

import (
	"reflect"
	"testing"

	"go.trulyao.dev/mirror/helper"
)

type Test struct {
	Description string
	Source      any
	Expected    Item
	WantErr     bool
}

func Test_ParseItem_Opts(t *testing.T) {
	type (
		OptTest struct {
			Description string
			Opt         Options
			Source      any
			Expected    Item
			WantErr     bool
		}

		Foo int
	)

	tests := []OptTest{
		{
			Description: "parse integer with nullable overridden to true",
			Opt:         Options{OverrideNullable: helper.Bool(true)},
			Source:      *new(Foo),
			Expected:    Scalar{"Foo", TypeInteger, true},
		},
	}

	for _, tt := range tests {
		got, err := ParseType(reflect.TypeOf(tt.Source), tt.Opt)
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

func Test_ParseItem_Scalar(t *testing.T) {
	type (
		// Ints
		Foo   int
		Foo8  int8
		Foo16 int16
		Foo32 int32
		Foo64 int64

		// Floats
		Float32 float32
		Float64 float64

		Language string

		IsEnabled bool
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
		{
			Description: "parse f32",
			Source:      *new(Float32),
			Expected:    Scalar{"Float32", TypeFloat, false},
		},
		{
			Description: "parse f64",
			Source:      *new(Float64),
			Expected:    Scalar{"Float64", TypeFloat, false},
		},
		{
			Description: "parse string",
			Source:      *new(Language),
			Expected:    Scalar{"Language", TypeString, false},
		},
		{
			Description: "parse boolean",
			Source:      *new(IsEnabled),
			Expected:    Scalar{"IsEnabled", TypeBoolean, false},
		},
	}

	for _, tt := range tests {
		got, err := ParseType(reflect.TypeOf(tt.Source))
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

func Test_ParseItem_Map(t *testing.T) {
	type (
		StringString map[string]string
		StringInt    map[string]int
		StringFloat  map[string]float32

		// TODO: test with struct as keys and values
		PtrStr      map[*string]*string
		ValuePtrStr map[string]*string
	)

	tests := []Test{
		{
			Description: "parse <string, string> map",
			Source:      StringString{},
			Expected: Map{
				"StringString",
				Scalar{"string", TypeString, false},
				Scalar{"string", TypeString, false},
			},
		},
		{
			Description: "parse <string, int> map",
			Source:      StringInt{},
			Expected: Map{
				"StringInt",
				Scalar{"string", TypeString, false},
				Scalar{"int", TypeInteger, false},
			},
		},
		{
			Description: "parse <string, float32> map",
			Source:      StringFloat{},
			Expected: Map{
				"StringFloat",
				Scalar{"string", TypeString, false},
				Scalar{"float32", TypeFloat, false},
			},
		},
		{
			Description: "parse <*string, *string> map",
			Source:      PtrStr{},
			Expected: Map{
				"PtrStr",
				Scalar{"string", TypeString, true},
				Scalar{"string", TypeString, true},
			},
		},
		{
			Description: "parse <string, *string> map",
			Source:      ValuePtrStr{},
			Expected: Map{
				"ValuePtrStr",
				Scalar{"string", TypeString, false},
				Scalar{"string", TypeString, true},
			},
		},
	}

	for _, tt := range tests {
		got, err := ParseType(reflect.TypeOf(tt.Source))
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
