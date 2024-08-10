package typescript_test

import (
	"testing"

	"go.trulyao.dev/mirror/generator/typescript"
	"go.trulyao.dev/mirror/parser"
)

type Test struct {
	Description string
	Config      typescript.Config
	Src         parser.Item
	Expect      string
	WantErr     bool
}

func Test_GenerateScalar(t *testing.T) {
	config := typescript.Config{
		InludeSemiColon: true,
	}

	tests := []Test{
		{
			Description: "generate string",
			Src: parser.Scalar{
				ItemName: "FooString",
				ItemType: parser.TypeString,
				Nullable: false,
			},
			Expect: "type FooString = string;",
			Config: config,
		},
		{
			Description: "generate nullable string",
			Src: parser.Scalar{
				ItemName: "NullableString",
				ItemType: parser.TypeString,
				Nullable: true,
			},
			Expect: "type NullableString = string | undefined;",
			Config: config,
		},
		{
			Description: "generate integer",
			Src: parser.Scalar{
				ItemName: "FooInt",
				ItemType: parser.TypeInteger,
				Nullable: false,
			},
			Expect: "type FooInt = number;",
			Config: config,
		},
		{
			Description: "generate nullable integer",
			Src: parser.Scalar{
				ItemName: "NullableInt",
				ItemType: parser.TypeInteger,
				Nullable: true,
			},
			Expect: "type NullableInt = number | undefined;",
			Config: config,
		},
	}

	runTests(t, tests)
}

func Test_GenerateArray(t *testing.T) {
	tests := []Test{
		{
			Description: "generate nullable integer array",
			Src: parser.List{
				ItemName: "IntArray",
				BaseItem: parser.Scalar{
					ItemName: "Int",
					ItemType: parser.TypeInteger,
					Nullable: false,
				},
				Nullable: true,
			},
			Expect: "type IntArray = Array<number> | null;",
			Config: typescript.Config{
				PreferArrayGeneric:    true,
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate integer|null array WITH generic array syntax",
			Src: parser.List{
				ItemName: "NullIntArray",
				BaseItem: parser.Scalar{
					ItemName: "Int",
					ItemType: parser.TypeInteger,
					Nullable: true,
				},
				Nullable: false,
			},
			Expect: "type NullIntArray = Array<number | null>;",
			Config: typescript.Config{
				PreferArrayGeneric:    true,
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate integer|null array WITHOUT generic array syntax",
			Src: parser.List{
				ItemName: "NullIntArray",
				BaseItem: parser.Scalar{
					ItemName: "Int",
					ItemType: parser.TypeInteger,
					Nullable: true,
				},
				Nullable: false,
			},
			Expect: "type NullIntArray = (number | null)[];",
			Config: typescript.Config{
				PreferArrayGeneric:    false,
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate object array",
			Src: parser.List{
				ItemName: "IntArray",
				BaseItem: parser.Struct{
					ItemName: "Foo",
					Fields:   []parser.Field{},
					Nullable: false,
				},
				Nullable: true,
			},
			Expect: "type IntArray = Array<Foo> | undefined;",
			Config: typescript.Config{
				PreferArrayGeneric: true,
				InludeSemiColon:    true,
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		gen := typescript.NewGenerator(&test.Config)
		got, err := gen.GenerateItem(test.Src)
		if err != nil {
			if !test.WantErr {
				t.Errorf("unexpected error: %v", err)
			}

			continue
		}

		if got != test.Expect {
			t.Errorf("expected %q, got %q", test.Expect, got)
		}
	}
}
