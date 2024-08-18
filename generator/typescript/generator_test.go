package typescript_test

import (
	"testing"

	"go.trulyao.dev/mirror/v2/config"
	"go.trulyao.dev/mirror/v2/generator/typescript"
	"go.trulyao.dev/mirror/v2/parser"
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
			Expect: "export type FooString = string;",
			Config: config,
		},
		{
			Description: "generate nullable string",
			Src: parser.Scalar{
				ItemName: "NullableString",
				ItemType: parser.TypeString,
				Nullable: true,
			},
			Expect: "export type NullableString = string | undefined;",
			Config: config,
		},
		{
			Description: "generate integer",
			Src: parser.Scalar{
				ItemName: "FooInt",
				ItemType: parser.TypeInteger,
				Nullable: false,
			},
			Expect: "export type FooInt = number;",
			Config: config,
		},
		{
			Description: "generate nullable integer",
			Src: parser.Scalar{
				ItemName: "NullableInt",
				ItemType: parser.TypeInteger,
				Nullable: true,
			},
			Expect: "export type NullableInt = number | undefined;",
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
			Expect: "export type IntArray = Array<number> | null;",
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
			Expect: "export type NullIntArray = Array<number | null>;",
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
			Expect: "export type NullIntArray = (number | null)[];",
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
			Expect: "export type IntArray = Array<Foo> | undefined;",
			Config: typescript.Config{
				PreferArrayGeneric: true,
				InludeSemiColon:    true,
			},
		},
	}

	runTests(t, tests)
}

func Test_GenerateStruct(t *testing.T) {
	tests := []Test{
		{
			Description: "generate struct",
			Src: parser.Struct{
				ItemName: "Foo",
				Fields: []parser.Field{
					{
						ItemName: "Bar",
						BaseItem: parser.Scalar{
							ItemName: "Baz",
							ItemType: parser.TypeString,
							Nullable: false,
						},
					},
				},
				Nullable: false,
			},
			Expect: "export type Foo = {\n    Bar: string;\n};",
			Config: typescript.Config{
				InludeSemiColon:  true,
				IndentationType:  config.IndentSpace,
				IndentationCount: 4,
				InlineObjects:    true,
			},
		},

		{
			Description: "generate struct with inline objects and string field (tab indentation)",
			Src: parser.Struct{
				ItemName: "Foo",
				Fields: []parser.Field{
					{
						ItemName: "Bar",
						BaseItem: parser.Scalar{
							ItemName: "Baz",
							ItemType: parser.TypeString,
							Nullable: false,
						},
					},
				},
				Nullable: false,
			},
			Expect: "export type Foo = {\n\tBar: string;\n};",
			Config: typescript.Config{
				InludeSemiColon:  true,
				IndentationType:  config.IndentTab,
				IndentationCount: 4,
				InlineObjects:    false,
			},
		},

		{
			Description: "generate struct with struct fields and inlining disabled",
			Src: parser.Struct{
				ItemName: "Foo",
				Fields: []parser.Field{
					{
						ItemName: "Bar",
						BaseItem: parser.Struct{
							ItemName: "Baz",
							Fields: []parser.Field{
								{
									ItemName: "Qux",
									BaseItem: parser.Scalar{
										ItemName: "Quux",
										ItemType: parser.TypeString,
									},
								},
							},
						},
					},
				},
			},
			Expect: "export type Foo = {\n\tBar: Baz;\n};",
			Config: typescript.Config{
				InludeSemiColon:  true,
				IndentationType:  config.IndentTab,
				IndentationCount: 4,
				InlineObjects:    false,
			},
		},

		{
			Description: "generate struct with struct fields and inlining ENABLED",
			Src: parser.Struct{
				ItemName: "Foo",
				Fields: []parser.Field{
					{
						ItemName: "Bar",
						BaseItem: parser.Struct{
							ItemName: "Baz",
							Fields: []parser.Field{
								{
									ItemName: "Qux",
									BaseItem: parser.Scalar{
										ItemName: "Quux",
										ItemType: parser.TypeInteger,
									},
								},
							},
						},
					},
				},
			},
			Expect: "export type Foo = {\n\tBar: {\n\t\tQux: number;\n\t};\n};",
			Config: typescript.Config{
				InludeSemiColon:  true,
				IndentationType:  config.IndentTab,
				IndentationCount: 4,
				InlineObjects:    true,
			},
		},

		{
			Description: "generate struct with array field (NO INLINING)",
			Src: parser.Struct{
				ItemName: "Foo",
				Fields: []parser.Field{
					{
						ItemName: "Bar",
						BaseItem: parser.List{
							ItemName: "Baz",
							BaseItem: parser.Scalar{
								ItemName: "",
								ItemType: parser.TypeBoolean,
							},
						},
					},
				},
			},
			Expect: "export type Foo = {\n\tBar: Array<boolean>;\n};",
			Config: typescript.Config{
				InludeSemiColon:    true,
				IndentationType:    config.IndentTab,
				IndentationCount:   4,
				InlineObjects:      false,
				PreferArrayGeneric: true,
			},
		},
	}

	runTests(t, tests)
}

func Test_GenerateMap(t *testing.T) {
	tests := []Test{
		{
			Description: "generate map with string key and integer value",
			Src: parser.Map{
				ItemName: "FooMap",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value:    parser.Scalar{ItemName: "integer", ItemType: parser.TypeInteger},
				Nullable: false,
			},
			Expect: "export type FooMap = Record<string, number>;",
			Config: typescript.Config{
				InludeSemiColon: true,
			},
		},

		{
			Description: "generate nullable map with string key and integer value",
			Src: parser.Map{
				ItemName: "FooMap",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value:    parser.Scalar{ItemName: "integer", ItemType: parser.TypeInteger},
				Nullable: true,
			},
			Expect: "export type FooMap = Record<string, number> | null;",
			Config: typescript.Config{
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate map with string key and nullable string value",
			Src: parser.Map{
				ItemName: "FooMapWithNullableValue",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value: parser.Scalar{
					ItemName: "string",
					ItemType: parser.TypeString,
					Nullable: true,
				},
				Nullable: false,
			},
			Expect: "export type FooMapWithNullableValue = Record<string, string | null>;",
			Config: typescript.Config{
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate map with non-scalar key",
			Src: parser.Map{
				ItemName: "FooMap",
				Key: parser.Struct{
					ItemName: "Foo",
					Fields:   []parser.Field{},
				},
				Value: parser.Scalar{ItemName: "integer", ItemType: parser.TypeInteger},
			},
			WantErr: true,
		},

		{
			Description: "generare map with nested map value (NO INLINING)",
			Src: parser.Map{
				ItemName: "FooMap",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value: parser.List{
					ItemName: "MapArray",
					BaseItem: parser.Map{
						ItemName: "NestedMap",
						Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
						Value:    parser.Scalar{ItemName: "integer", ItemType: parser.TypeInteger},
					},
				},
				Nullable: false,
			},

			Expect: "export type FooMap = Record<string, Array<NestedMap>>;",
			Config: typescript.Config{
				InludeSemiColon:    true,
				PreferArrayGeneric: true,
			},
		},

		{
			Description: "generare map with nested map value (INLINING ENABLED)",
			Src: parser.Map{
				ItemName: "FooMap",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value: parser.List{
					ItemName: "MapArray",
					BaseItem: parser.Map{
						ItemName: "NestedMap",
						Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
						Value:    parser.Scalar{ItemName: "integer", ItemType: parser.TypeInteger},
					},
				},
				Nullable: false,
			},

			Expect: "export type FooMap = Record<string, Array<Record<string, number>>>;",
			Config: typescript.Config{
				InludeSemiColon:    true,
				PreferArrayGeneric: true,
				InlineObjects:      true,
			},
		},

		{
			Description: "generare map with two nested maps (INLINING ENABLED)",
			Src: parser.Map{
				ItemName: "FooMap",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value: parser.List{
					ItemName: "MapArray",
					BaseItem: parser.Map{
						ItemName: "NestedMap",
						Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
						Value: parser.Map{
							ItemName: "InnerNestedMap",
							Key: parser.Scalar{
								ItemName: "string",
								ItemType: parser.TypeString,
							},
							Value: parser.Scalar{
								ItemName: "integer",
								ItemType: parser.TypeInteger,
							},
						},
					},
				},
				Nullable: false,
			},

			Expect: "export type FooMap = Record<string, Array<Record<string, Record<string, number>>>>;",
			Config: typescript.Config{
				InludeSemiColon:    true,
				PreferArrayGeneric: true,
				InlineObjects:      true,
			},
		},

		{
			Description: "generare map with two nested maps (NO INLINING)",
			Src: parser.Map{
				ItemName: "FooMap",
				Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				Value: parser.List{
					ItemName: "MapArray",
					BaseItem: parser.Map{
						ItemName: "NestedMap",
						Key:      parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
						Value: parser.Map{
							ItemName: "InnerNestedMap",
							Key: parser.Scalar{
								ItemName: "string",
								ItemType: parser.TypeString,
							},
							Value: parser.Scalar{
								ItemName: "integer",
								ItemType: parser.TypeInteger,
							},
						},
					},
				},
				Nullable: false,
			},

			Expect: "export type FooMap = Record<string, Array<NestedMap>>;",
			Config: typescript.Config{
				InludeSemiColon:    true,
				PreferArrayGeneric: true,
				InlineObjects:      false,
			},
		},
	}

	runTests(t, tests)
}

func Test_GenerateFunc(t *testing.T) {
	tests := []Test{
		{
			Description: "generate function with no params or returns",
			Src: parser.Function{
				ItemName: "VoidFunc",
				Params:   []parser.Item{},
				Returns:  []parser.Item{},
				Nullable: false,
			},
			Expect: "export type VoidFunc = () => void;",
			Config: typescript.Config{
				InludeSemiColon: true,
			},
		},

		{
			Description: "generate function with no params and single returns",
			Src: parser.Function{
				ItemName: "SingleReturnFunc",
				Params:   []parser.Item{},
				Returns: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				},
				Nullable: false,
			},
			Expect: "export type SingleReturnFunc = () => string;",
			Config: typescript.Config{
				InludeSemiColon: true,
			},
		},

		{
			Description: "generate function with multiple params and returns",
			Src: parser.Function{
				ItemName: "MultiFunc",
				Params: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
					parser.Scalar{ItemName: "number", ItemType: parser.TypeInteger},
				},
				Returns: []parser.Item{
					parser.Scalar{ItemName: "boolean", ItemType: parser.TypeBoolean},
				},
			},
			Expect: "export type MultiFunc = (arg0: string, arg1: number) => boolean;",
			Config: typescript.Config{InludeSemiColon: true},
		},

		{
			Description: "generate function with multiple params and multiple returns",
			Src: parser.Function{
				ItemName: "MultiFunc",
				Params: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
					parser.Scalar{ItemName: "number", ItemType: parser.TypeInteger},
				},
				Returns: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
					parser.Scalar{ItemName: "boolean", ItemType: parser.TypeBoolean},
				},
			},
			Expect:  "",
			WantErr: true,
			Config:  typescript.Config{InludeSemiColon: true},
		},

		{
			Description: "generate function with nullable return",
			Src: parser.Function{
				ItemName: "NullableReturnFunc",
				Params:   []parser.Item{},
				Returns: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString, Nullable: true},
				},
			},
			Expect: "export type NullableReturnFunc = () => string | null;",
			Config: typescript.Config{
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate function with nullable params",
			Src: parser.Function{
				ItemName: "NullableParamFunc",
				Params: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString, Nullable: true},
					parser.Scalar{
						ItemName: "number",
						ItemType: parser.TypeInteger,
						Nullable: false,
					},
				},
				Returns: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
				},
			},
			Expect: "export type NullableParamFunc = (arg0: string | null, arg1: number) => string;",
			Config: typescript.Config{
				InludeSemiColon:       true,
				PreferNullForNullable: true,
			},
		},

		{
			Description: "generate function with function param (INLINING ENABLED)",
			Src: parser.Function{
				ItemName: "FuncParamFunc",
				Params: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
					parser.Function{
						ItemName: "InnerFunc",
						Params:   []parser.Item{},
						Returns: []parser.Item{
							parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
						},
					},
				},
				Returns: []parser.Item{},
			},
			Expect: "export type FuncParamFunc = (arg0: string, arg1: () => string) => void;",
			Config: typescript.Config{
				InlineObjects:   true,
				InludeSemiColon: true,
			},
		},

		{
			Description: "generate function with function param (NO INLINING)",
			Src: parser.Function{
				ItemName: "FuncParamFunc",
				Params: []parser.Item{
					parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
					parser.Function{
						ItemName: "InnerFunc",
						Params:   []parser.Item{},
						Returns: []parser.Item{
							parser.Scalar{ItemName: "string", ItemType: parser.TypeString},
						},
					},
				},
				Returns: []parser.Item{},
			},
			Expect: "export type FuncParamFunc = (arg0: string, arg1: InnerFunc) => void;",
			Config: typescript.Config{
				InlineObjects:   false,
				InludeSemiColon: true,
			},
		},

		{
			Description: "generate function with function return",
			Src: parser.Function{
				ItemName: "FuncReturnFunc",
				Params:   []parser.Item{},
				Returns: []parser.Item{
					parser.Function{
						ItemName: "InnerFunc",
						Params:   []parser.Item{},
						Returns:  []parser.Item{},
					},
				},
			},
			Expect: "export type FuncReturnFunc = () => (() => void);",
			Config: typescript.Config{
				InludeSemiColon: true,
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		gen := typescript.NewGenerator(&test.Config)
		gen.SetNonStrict(true)

		got, err := gen.GenerateItem(test.Src)
		if err != nil {
			if !test.WantErr {
				t.Errorf("[%s] unexpected error: %v", test.Description, err)
			}

			continue
		}

		if got != test.Expect {
			t.Errorf("[%s] expected %q, got %q", test.Description, test.Expect, got)
		}
	}
}
