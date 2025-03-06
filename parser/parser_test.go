package parser_test

import (
	"database/sql"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"go.trulyao.dev/mirror/v2/extractor/meta"
	"go.trulyao.dev/mirror/v2/parser"
)

// TODO: add tests for a mix of JSON and MIRROR tags

type Test struct {
	Description string
	Source      any
	Expected    parser.Item
	WantErr     bool
}

func Test_ParseItem_Opts(t *testing.T) {
	type (
		OptTest struct {
			Description string
			Opt         parser.Options
			Source      any
			Expected    parser.Item
			WantErr     bool
		}

		Foo int
	)

	tests := []OptTest{
		{
			Description: "parse integer with nullable overridden to true",
			Opt:         parser.Options{OverrideNullable: true},
			Source:      *new(Foo),
			Expected:    &parser.Scalar{"Foo", parser.TypeInteger, true},
		},
	}

	for _, tt := range tests {
		p := parser.New()

		got, err := p.ParseWithOpts(reflect.TypeOf(tt.Source), tt.Opt)
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
			Expected:    &parser.Scalar{"Foo", parser.TypeInteger, false},
		},
		{
			Description: "parse i8",
			Source:      *new(Foo8),
			Expected:    &parser.Scalar{"Foo8", parser.TypeInteger, false},
		},
		{
			Description: "parse i16",
			Source:      *new(Foo16),
			Expected:    &parser.Scalar{"Foo16", parser.TypeInteger, false},
		},
		{
			Description: "parse i32",
			Source:      *new(Foo32),
			Expected:    &parser.Scalar{"Foo32", parser.TypeInteger, false},
		},
		{
			Description: "parse i64",
			Source:      *new(Foo64),
			Expected:    &parser.Scalar{"Foo64", parser.TypeInteger, false},
		},
		{
			Description: "parse f32",
			Source:      *new(Float32),
			Expected:    &parser.Scalar{"Float32", parser.TypeFloat, false},
		},
		{
			Description: "parse f64",
			Source:      *new(Float64),
			Expected:    &parser.Scalar{"Float64", parser.TypeFloat, false},
		},
		{
			Description: "parse string",
			Source:      *new(Language),
			Expected:    &parser.Scalar{"Language", parser.TypeString, false},
		},
		{
			Description: "parse boolean",
			Source:      *new(IsEnabled),
			Expected:    &parser.Scalar{"IsEnabled", parser.TypeBoolean, false},
		},
	}

	runTests(t, tests)
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
			Expected: &parser.Map{
				"StringString",
				&parser.Scalar{"string", parser.TypeString, false},
				&parser.Scalar{"string", parser.TypeString, false},
				false,
			},
		},
		{
			Description: "parse <string, int> map",
			Source:      StringInt{},
			Expected: &parser.Map{
				"StringInt",
				&parser.Scalar{"string", parser.TypeString, false},
				&parser.Scalar{"int", parser.TypeInteger, false},
				false,
			},
		},
		{
			Description: "parse <string, float32> map",
			Source:      StringFloat{},
			Expected: &parser.Map{
				"StringFloat",
				&parser.Scalar{"string", parser.TypeString, false},
				&parser.Scalar{"float32", parser.TypeFloat, false},
				false,
			},
		},
		{
			Description: "parse <*string, *string> map",
			Source:      PtrStr{},
			Expected: &parser.Map{
				"PtrStr",
				&parser.Scalar{"string", parser.TypeString, true},
				&parser.Scalar{"string", parser.TypeString, true},
				false,
			},
		},
		{
			Description: "parse <string, *string> map",
			Source:      ValuePtrStr{},
			Expected: &parser.Map{
				"ValuePtrStr",
				&parser.Scalar{"string", parser.TypeString, false},
				&parser.Scalar{"string", parser.TypeString, true},
				false,
			},
		},
	}

	runTests(t, tests)
}

func Test_ParseItem_Struct(t *testing.T) {
	type Person struct {
		FirstName string
		LastName  string
	}

	type User struct {
		FullName *string `json:"name"  mirror:"name:full_name,optional:false"`
		Username string  `json:"uname"`
		Password string  `json:"pass"`
	}

	type Account struct {
		User      *User `mirror:"name:linked_user"`
		CreatedAt int   `mirror:"name:created_at"`
	}

	type Meta struct {
		ID        int32     `json:"-"          mirror:"type:string"`
		Scope     string    `json:"scope"      mirror:"type:'reset' | 'change'"`
		CreatedAt time.Time `json:"created_at" mirror:"type:Date,skip:true,optional:true"`
	}

	tests := []Test{
		{
			Description: "parse Person struct",
			Source:      Person{},
			Expected: &parser.Struct{
				"Person",
				[]parser.Field{
					{
						ItemName: "FirstName",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "FirstName",
							Name:         "FirstName",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					{
						ItemName: "LastName",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "LastName",
							Name:         "LastName",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
				false,
			},
		},

		{
			Description: "parse User struct with JSON meta",
			Source:      User{},
			Expected: &parser.Struct{
				ItemName: "User",
				Fields: []parser.Field{
					{
						ItemName: "full_name",
						BaseItem: &parser.Scalar{"string", parser.TypeString, true},
						Meta: meta.Meta{
							OriginalName: "FullName",
							Name:         "full_name",
							Type:         "",
							Optional:     meta.OptionalFalse,
							Skip:         false,
						},
					},

					{
						ItemName: "uname",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Username",
							Name:         "uname",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					{
						ItemName: "pass",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Password",
							Name:         "pass",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
				Nullable: false,
			},
		},

		{
			Description: "parse Account struct with mirror meta and pointer field",
			Source:      Account{},
			Expected: &parser.Struct{
				ItemName: "Account",
				Fields: []parser.Field{
					{
						ItemName: "linked_user",
						BaseItem: &parser.Struct{
							ItemName: "User",
							Fields: []parser.Field{
								{
									ItemName: "full_name",
									BaseItem: &parser.Scalar{"string", parser.TypeString, true},
									Meta: meta.Meta{
										OriginalName: "FullName",
										Name:         "full_name",
										Type:         "",
										Optional:     meta.OptionalFalse,
										Skip:         false,
									},
								},
								{
									ItemName: "uname",
									BaseItem: &parser.Scalar{"string", parser.TypeString, false},
									Meta: meta.Meta{
										OriginalName: "Username",
										Name:         "uname",
										Type:         "",
										Optional:     meta.OptionalNone,
										Skip:         false,
									},
								},

								{
									ItemName: "pass",
									BaseItem: &parser.Scalar{"string", parser.TypeString, false},
									Meta: meta.Meta{
										OriginalName: "Password",
										Name:         "pass",
										Type:         "",
										Optional:     meta.OptionalNone,
										Skip:         false,
									},
								},
							},
							Nullable: true,
						},
						Meta: meta.Meta{
							OriginalName: "User",
							Name:         "linked_user",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					{
						ItemName: "created_at",
						BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
						Meta: meta.Meta{
							OriginalName: "CreatedAt",
							Name:         "created_at",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
				Nullable: false,
			},
		},

		{
			Description: "parse Meta struct with mixed JSON and mirror meta",
			Source:      Meta{},
			Expected: &parser.Struct{
				ItemName: "Meta",
				Fields: []parser.Field{
					{
						ItemName: "ID",
						BaseItem: &parser.Scalar{"int32", parser.TypeInteger, false},
						Meta: meta.Meta{
							OriginalName: "ID",
							Name:         "ID",
							Type:         "string",
							Optional:     meta.OptionalNone,
							Skip:         true,
						},
					},
					{
						ItemName: "scope",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Scope",
							Name:         "scope",
							Type:         "'reset' | 'change'",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
					{
						ItemName: "created_at",
						BaseItem: &parser.Scalar{"Time", parser.TypeTimestamp, false},
						Meta: meta.Meta{
							OriginalName: "CreatedAt",
							Name:         "created_at",
							Type:         "Date",
							Optional:     meta.OptionalTrue,
							Skip:         true,
						},
					},
				},
			},
		},
	}

	runTests(t, tests)
}

func Test_ParseItem_List(t *testing.T) {
	type (
		CustomType struct{ Name string }

		Strings    []string
		Ints       []int
		Floats     []float32
		Structs    []CustomType
		StringPtrs []*string
		ListList   [][]int
		ListPtr    *[]int
		ListPtrs   []*[]int

		FixedStrings [3]string
		FixedStructs [8]CustomType
		FixedIntPtrs [6]*int
	)

	tests := []Test{
		// Strings
		{
			Description: "parse []string",
			Source:      Strings{},
			Expected: &parser.List{
				ItemName: "Strings",
				BaseItem: &parser.Scalar{"string", parser.TypeString, false},
				Nullable: false,
				Length:   parser.EmptyLength,
			},
		},

		// Ints
		{
			Description: "parse []ints",
			Source:      Ints{},
			Expected: &parser.List{
				ItemName: "Ints",
				BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
				Nullable: false,
				Length:   parser.EmptyLength,
			},
		},

		// Floats
		{
			Description: "parse []floats",
			Source:      Floats{},
			Expected: &parser.List{
				ItemName: "Floats",
				BaseItem: &parser.Scalar{"float32", parser.TypeFloat, false},
				Nullable: false,
				Length:   parser.EmptyLength,
			},
		},

		// Structs
		{
			Description: "parse []structs",
			Source:      Structs{},
			Expected: &parser.List{
				ItemName: "Structs",
				BaseItem: &parser.Struct{
					ItemName: "CustomType",
					Fields: []parser.Field{
						{
							ItemName: "Name",
							BaseItem: &parser.Scalar{"string", parser.TypeString, false},
							Meta: meta.Meta{
								OriginalName: "Name",
								Name:         "Name",
								Type:         "",
								Optional:     meta.OptionalNone,
								Skip:         false,
							},
						},
					},
				},
				Nullable: false,
				Length:   parser.EmptyLength,
			},
		},

		// String pointers
		{
			Description: "parse []*string",
			Source:      StringPtrs{},
			Expected: &parser.List{
				ItemName: "StringPtrs",
				BaseItem: &parser.Scalar{"string", parser.TypeString, true},
				Nullable: false,
				Length:   parser.EmptyLength,
			},
		},

		// List of lists
		{
			Description: "parse [][]int",
			Source:      ListList{},
			Expected: &parser.List{
				ItemName: "ListList",
				BaseItem: &parser.List{
					ItemName: "", // The inner list has no name
					BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
					Length:   parser.EmptyLength,
					Nullable: false,
				},
				Nullable: false,
				Length:   parser.EmptyLength,
			},
		},

		// List pointer
		{
			Description: "parse *[]int",
			Source:      *new(ListPtr), // new(ListPtr) returns a pointer to a nil slice, that is intentionally unhandled by the parser and will return an error for now
			Expected: &parser.List{
				ItemName: "",
				BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
				Nullable: true,
				Length:   parser.EmptyLength,
			},
		},

		// List of list pointers
		{
			Description: "parse []*[]int",
			Source:      ListPtrs{},
			Expected: &parser.List{
				ItemName: "ListPtrs",
				Length:   parser.EmptyLength,
				BaseItem: &parser.List{
					ItemName: "",
					BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
					Length:   parser.EmptyLength,
					Nullable: true,
				},
			},
		},

		// Fixed strings
		{
			Description: "parse [3]string",
			Source:      FixedStrings{},
			Expected: &parser.List{
				ItemName: "FixedStrings",
				BaseItem: &parser.Scalar{"string", parser.TypeString, false},
				Length:   3,
				Nullable: false,
			},
		},

		// Fixed Structs
		{
			Description: "parse [8]structs",
			Source:      FixedStructs{},
			Expected: &parser.List{
				ItemName: "FixedStructs",
				BaseItem: &parser.Struct{
					ItemName: "CustomType",
					Fields: []parser.Field{
						{
							ItemName: "Name",
							BaseItem: &parser.Scalar{"string", parser.TypeString, false},
							Meta: meta.Meta{
								OriginalName: "Name",
								Name:         "Name",
								Type:         "",
								Optional:     meta.OptionalNone,
								Skip:         false,
							},
						},
					},
				},
				Nullable: false,
				Length:   8,
			},
		},

		// Fixed IntPtrs
		{
			Description: "parse [6]*int",
			Source:      FixedIntPtrs{},
			Expected: &parser.List{
				ItemName: "FixedIntPtrs",
				BaseItem: &parser.Scalar{"int", parser.TypeInteger, true},
				Length:   6,
				Nullable: false,
			},
		},
	}

	runTests(t, tests)
}

func Test_ParseItem_Function(t *testing.T) {
	type (
		// Function types
		Func1          func() error
		Add            func(int, int) int
		ReturnMultiple func(string, *string) (int, error)

		Foo       struct{ Name string }
		InsertFoo func(Foo) error
	)

	tests := []Test{
		{
			Description: "parse func() error",
			Source:      Func1(nil),
			Expected: &parser.Function{
				ItemName: "Func1",
				Params:   []parser.Item{},
				Returns:  []parser.Item{&parser.Scalar{"error", parser.TypeString, false}},
				Nullable: false,
			},
		},

		// Add
		{
			Description: "parse func(int, int) int",
			Source:      Add(nil),
			Expected: &parser.Function{
				ItemName: "Add",
				Params: []parser.Item{
					&parser.Scalar{"int", parser.TypeInteger, false},
					&parser.Scalar{"int", parser.TypeInteger, false},
				},
				Returns:  []parser.Item{&parser.Scalar{"int", parser.TypeInteger, false}},
				Nullable: false,
			},
		},

		// ReturnMultiple
		{
			Description: "parse func(string, string) (int, error)",
			Source:      ReturnMultiple(nil),
			Expected: &parser.Function{
				ItemName: "ReturnMultiple",
				Params: []parser.Item{
					&parser.Scalar{"string", parser.TypeString, false},
					&parser.Scalar{"string", parser.TypeString, true},
				},
				Returns: []parser.Item{
					&parser.Scalar{"int", parser.TypeInteger, false},
					&parser.Scalar{"error", parser.TypeString, false},
				},
				Nullable: false,
			},
		},

		// InsertFoo
		{
			Description: "parse func(Foo) error",
			Source:      InsertFoo(nil),
			Expected: &parser.Function{
				ItemName: "InsertFoo",
				Params: []parser.Item{
					&parser.Struct{
						ItemName: "Foo",
						Fields: []parser.Field{
							{
								ItemName: "Name",
								BaseItem: &parser.Scalar{"string", parser.TypeString, false},
								Meta: meta.Meta{
									OriginalName: "Name",
									Name:         "Name",
									Type:         "",
									Optional:     meta.OptionalNone,
									Skip:         false,
								},
							},
						},
						Nullable: false,
					},
				},
				Returns:  []parser.Item{&parser.Scalar{"error", parser.TypeString, false}},
				Nullable: false,
			},
		},
	}

	runTests(t, tests)
}

func Test_ParseEmbeddedStruct(t *testing.T) {
	type (
		FooEmbedded struct{ Name string }
		FooParent   struct {
			FooEmbedded
			Age int
		}

		EmbeddedString        string
		EmbeddedInt           int
		EmbeddedBool          bool
		FooWithEmbeddedString struct {
			EmbeddedString `mirror:"name:embedded_string"` // with custom tag
			*EmbeddedInt                                   // without custom tag and optional
			*EmbeddedBool  `mirror:"name:probably,type:number"`
		}
	)

	tests := []Test{
		{
			Description: "parse embedded struct",
			Source:      FooParent{},
			Expected: &parser.Struct{
				ItemName: "FooParent",
				Fields: []parser.Field{
					{
						ItemName: "Name",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Name",
							Name:         "Name",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
					{
						ItemName: "Age",
						BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
						Meta: meta.Meta{
							OriginalName: "Age",
							Name:         "Age",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
			},
		},

		{
			Description: "parse struct with embedded non-struct type",
			Source:      FooWithEmbeddedString{},
			Expected: &parser.Struct{
				ItemName: "FooWithEmbeddedString",
				Fields: []parser.Field{
					{
						ItemName: "embedded_string",
						BaseItem: &parser.Scalar{"EmbeddedString", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "EmbeddedString",
							Name:         "embedded_string",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					{
						ItemName: "EmbeddedInt",
						BaseItem: &parser.Scalar{"EmbeddedInt", parser.TypeInteger, true},
						Meta: meta.Meta{
							OriginalName: "EmbeddedInt",
							Name:         "EmbeddedInt",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					{
						ItemName: "probably",
						BaseItem: &parser.Scalar{"EmbeddedBool", parser.TypeBoolean, true},
						Meta: meta.Meta{
							OriginalName: "EmbeddedBool",
							Name:         "probably",
							Type:         "number",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
			},
		},
	}

	runTests(t, tests)
}

func Test_ParseBuiltInTypes(t *testing.T) {
	type (
		TimeSlice []time.Time
		TimeArray [3]time.Time
	)

	tests := []Test{
		{
			Description: "parse time.Time",
			Source:      time.Time{},
			Expected:    &parser.Scalar{"Time", parser.TypeTimestamp, false},
		},

		{
			Description: "parse nullable time.Time",
			Source:      &time.Time{},
			Expected:    &parser.Scalar{"Time", parser.TypeTimestamp, true},
		},

		{
			Description: "parse []time.Time",
			Source:      TimeSlice{},
			Expected: &parser.List{
				"TimeSlice",
				&parser.Scalar{"Time", parser.TypeTimestamp, false},
				false,
				parser.EmptyLength,
			},
		},

		{
			Description: "parse []time.Time",
			Source:      &TimeArray{},
			Expected: &parser.List{
				"TimeArray",
				&parser.Scalar{"Time", parser.TypeTimestamp, false},
				true,
				3,
			},
		},

		{
			Description: "parse time.Duration",
			Source:      time.Duration(0),
			Expected:    &parser.Scalar{"Duration", parser.TypeInteger, false},
		},

		{
			Description: "parse nullable time.Duration",
			Source:      new(time.Duration),
			Expected:    &parser.Scalar{"Duration", parser.TypeInteger, true},
		},

		{
			Description: "parse sql.NullTime",
			Source:      sql.NullTime{},
			Expected:    &parser.Scalar{"NullTime", parser.TypeTimestamp, true},
		},

		{
			Description: "parse sql.NullInt64",
			Source:      sql.NullInt64{},
			Expected:    &parser.Scalar{"NullInt64", parser.TypeInteger, true},
		},
	}

	runTests(t, tests)
}

func Test_ParserHooks(t *testing.T) {
	type (
		TargetFoo struct {
			Name string
		}

		NotTargetFoo struct {
			Name string
		}

		Person struct {
			FirstName string
			LastName  string
		}
	)

	tests := []Test{
		{
			Description: "parse struct with OnParseItem hook, add a new field",
			Source:      TargetFoo{},
			Expected: &parser.Struct{
				ItemName: "TargetFoo",
				Fields: []parser.Field{
					{
						ItemName: "Name",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Name",
							Name:         "Name",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					// Added dynamically
					{
						ItemName: "Age",
						BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
						Meta: meta.Meta{
							OriginalName: "Age",
							Name:         "AddedAge",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
			},
		},

		{
			Description: "parse struct with OnParseItem hook, do not add a new field",
			Source:      NotTargetFoo{},
			Expected: &parser.Struct{
				ItemName: "NotTargetFoo",
				Fields: []parser.Field{
					{
						ItemName: "Name",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Name",
							Name:         "Name",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
			},
		},

		{
			Description: "parse struct with OnParseField hook, renaming FirstName to FName",
			Source:      Person{},
			Expected: &parser.Struct{
				ItemName: "Person",
				Fields: []parser.Field{
					{
						ItemName: "FName",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "FirstName",
							Name:         "FName",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},

					{
						ItemName: "LastName",
						BaseItem: &parser.Scalar{"string", parser.TypeString, false},
						Meta: meta.Meta{
							OriginalName: "LastName",
							Name:         "LastName",
							Type:         "",
							Optional:     meta.OptionalNone,
							Skip:         false,
						},
					},
				},
			},
		},
	}

	p := parser.New()
	p.SetFlattenEmbeddedTypes(true)

	p.OnParseItem(func(name string, item parser.Item) error {
		// Add a new field to the struct if the name is "TargetFoo"
		if name == "TargetFoo" {
			item.(*parser.Struct).Fields = append(item.(*parser.Struct).Fields, parser.Field{
				ItemName: "Age",
				BaseItem: &parser.Scalar{"int", parser.TypeInteger, false},
				Meta: meta.Meta{
					OriginalName: "Age",
					Name:         "AddedAge",
					Type:         "",
					Optional:     meta.OptionalNone,
					Skip:         false,
				},
			})
		}

		return nil
	})

	p.OnParseField(
		func(parentType *reflect.Type, originalField *reflect.StructField, field *parser.Field) error {
			// Modify the field if the original name is "FirstName"
			if field.Meta.OriginalName == "FirstName" {
				field.Meta.Name = "FName"
				field.ItemName = "FName"
			}

			return nil
		},
	)

	for _, tt := range tests {
		got, err := p.Parse(reflect.TypeOf(tt.Source))
		if err != nil && !tt.WantErr {
			t.Errorf("[%s] wanted NO error, got error `%s`", tt.Description, err.Error())
		}

		if err == nil && tt.WantErr {
			t.Errorf("[%s] wanted error, got no error", tt.Description)
		}

		if !reflect.DeepEqual(got, tt.Expected) {
			t.Errorf("[%s] wanted %#v, got %#v", tt.Description, tt.Expected, got)
		}
	}
}

func Test_StrucMethods(t *testing.T) {
	type (
		Person struct {
			FirstName string `mirror:"name:first_name"`
			LastName  string
		}
	)

	p := parser.New()
	parsed, err := p.Parse(reflect.TypeOf(Person{}))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var (
		parsedStruct *parser.Struct
		ok           bool
	)
	if parsedStruct, ok = parsed.(*parser.Struct); !ok {
		t.Fatalf("unexpected type: %T", parsed)
	}

	fname, found := parsedStruct.GetField("first_name")
	if !found {
		t.Fatalf("field not found: first_name")
	}

	if fname.Meta.OriginalName != "FirstName" {
		t.Errorf("expected original name to be FirstName, got %s", fname.Meta.OriginalName)
	}

	lname, found := parsedStruct.GetField("LastName")
	if !found {
		t.Fatalf("field not found: LastName")
	}

	if lname.Meta.OriginalName != "LastName" {
		t.Errorf("expected original name to be LastName, got %s", lname.Meta.OriginalName)
	}

	fnameIndex := parsedStruct.GetFieldIndex("first_name")
	if fnameIndex == -1 {
		t.Fatalf("field index not found but expected: first_name")
	}

	if fnameIndex != 0 {
		t.Errorf("expected field index to be 0, got %d", fnameIndex)
	}

	lnameIndex := parsedStruct.GetFieldIndex("LastName")
	if lnameIndex == -1 {
		t.Fatalf("field index not found but expected: LastName")
	}

	if lnameIndex != 1 {
		t.Errorf("expected field index to be 1, got %d", lnameIndex)
	}

	_, found = parsedStruct.GetFieldByOriginalName("first_name")
	if found {
		t.Fatalf("expected field not to be found: first_name")
	}

	lnameByOrig, found := parsedStruct.GetFieldByOriginalName("LastName")
	if !found {
		t.Fatalf("field not found: LastName")
	}

	if lnameByOrig.Meta.OriginalName != "LastName" {
		t.Errorf("expected original name to be LastName, got %s", lnameByOrig.Meta.OriginalName)
	}
}

func Test_ParseCustomItem(t *testing.T) {
	type (
		__internal_scalar_type string
		__internal_struct_type struct {
			Name string
		}

		__internal_unregistered_type string
	)

	var (
		internalScalarItem = &parser.Scalar{"overriden_scalar_type", parser.TypeVoid, false}

		internalStructItem = &parser.Struct{
			ItemName: "overriden_struct_type",
			Fields: []parser.Field{
				{
					ItemName: "name",
					BaseItem: &parser.Scalar{"string", parser.TypeString, false},
					Meta:     meta.Meta{},
				},
			},
		}
	)

	tests := []Test{
		{
			Description: "parse custom scalar item",
			Source:      __internal_scalar_type(""),
			Expected:    internalScalarItem,
		},

		{
			Description: "parse custom struct item",
			Source:      __internal_struct_type{},
			Expected:    internalStructItem,
		},

		{
			Description: "parse unregistered custom type",
			Source:      __internal_unregistered_type(""),
			Expected:    &parser.Scalar{"__internal_unregistered_type", parser.TypeString, false},
		},
	}

	p := parser.New()
	p.AddCustomTypes([]parser.CustomType{
		{"__internal_scalar_type", internalScalarItem},
		{"__internal_struct_type", internalStructItem},
	})

	runTests(t, tests, p)
}

func runTests(t *testing.T, tests []Test, optParser ...*parser.Parser) {
	for _, tt := range tests {
		runTest(t, tt, optParser...)
	}
}

func runTest(t *testing.T, tt Test, optParse ...*parser.Parser) {
	var p *parser.Parser

	if len(optParse) > 0 {
		p = optParse[0]
	} else {
		p = parser.New()
		p.SetFlattenEmbeddedTypes(true) // TODO: make this better
	}

	got, err := p.Parse(reflect.TypeOf(tt.Source))
	if err != nil && !tt.WantErr {
		t.Errorf("[%s] wanted NO error, got error `%s`", tt.Description, err.Error())
	}

	if err == nil && tt.WantErr {
		t.Errorf("[%s] wanted error, got no error", tt.Description)
	}

	if !reflect.DeepEqual(got, tt.Expected) {
		if os.Getenv("JSON_DEBUG") == "true" {
			var (
				gotJson      []byte
				expectedJson []byte
			)

			gotJson, _ = json.MarshalIndent(got, "", "  ")
			expectedJson, _ = json.MarshalIndent(tt.Expected, "", "  ")

			t.Logf("got:\n%s", gotJson)
			t.Logf("expected:\n%s", expectedJson)
		}

		t.Errorf("[%s] wanted %#v, got %#v", tt.Description, tt.Expected, got)
	}
}
