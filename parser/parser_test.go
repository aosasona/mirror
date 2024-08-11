package parser

import (
	"reflect"
	"testing"

	"go.trulyao.dev/mirror/extractor/meta"
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
			Opt:         Options{OverrideNullable: true},
			Source:      *new(Foo),
			Expected:    Scalar{"Foo", TypeInteger, true},
		},
	}

	for _, tt := range tests {
		p := New()

		got, err := p.Parse(reflect.TypeOf(tt.Source), tt.Opt)
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
			Expected: Map{
				"StringString",
				Scalar{"string", TypeString, false},
				Scalar{"string", TypeString, false},
				false,
			},
		},
		{
			Description: "parse <string, int> map",
			Source:      StringInt{},
			Expected: Map{
				"StringInt",
				Scalar{"string", TypeString, false},
				Scalar{"int", TypeInteger, false},
				false,
			},
		},
		{
			Description: "parse <string, float32> map",
			Source:      StringFloat{},
			Expected: Map{
				"StringFloat",
				Scalar{"string", TypeString, false},
				Scalar{"float32", TypeFloat, false},
				false,
			},
		},
		{
			Description: "parse <*string, *string> map",
			Source:      PtrStr{},
			Expected: Map{
				"PtrStr",
				Scalar{"string", TypeString, true},
				Scalar{"string", TypeString, true},
				false,
			},
		},
		{
			Description: "parse <string, *string> map",
			Source:      ValuePtrStr{},
			Expected: Map{
				"ValuePtrStr",
				Scalar{"string", TypeString, false},
				Scalar{"string", TypeString, true},
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
		Username string `json:"uname"`
		Password string `json:"pass"`
	}

	type Account struct {
		User      *User `mirror:"name:linked_user"`
		CreatedAt int   `mirror:"name:created_at"`
	}

	tests := []Test{
		{
			Description: "parse Person struct",
			Source:      Person{},
			Expected: Struct{
				"Person",
				[]Field{
					{
						ItemName: "FirstName",
						BaseItem: Scalar{"string", TypeString, false},
						Meta: meta.Meta{
							OriginalName: "FirstName",
							Name:         "FirstName",
							Type:         "",
							Optional:     false,
							Skip:         false,
						},
					},

					{
						ItemName: "LastName",
						BaseItem: Scalar{"string", TypeString, false},
						Meta: meta.Meta{
							OriginalName: "LastName",
							Name:         "LastName",
							Type:         "",
							Optional:     false,
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
			Expected: Struct{
				ItemName: "User",
				Fields: []Field{
					{
						ItemName: "uname",
						BaseItem: Scalar{"string", TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Username",
							Name:         "uname",
							Type:         "",
							Optional:     false,
							Skip:         false,
						},
					},

					{
						ItemName: "pass",
						BaseItem: Scalar{"string", TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Password",
							Name:         "pass",
							Type:         "",
							Optional:     false,
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
			Expected: Struct{
				ItemName: "Account",
				Fields: []Field{
					{
						ItemName: "linked_user",
						BaseItem: Struct{
							ItemName: "User",
							Fields: []Field{
								{
									ItemName: "uname",
									BaseItem: Scalar{"string", TypeString, false},
									Meta: meta.Meta{
										OriginalName: "Username",
										Name:         "uname",
										Type:         "",
										Optional:     false,
										Skip:         false,
									},
								},

								{
									ItemName: "pass",
									BaseItem: Scalar{"string", TypeString, false},
									Meta: meta.Meta{
										OriginalName: "Password",
										Name:         "pass",
										Type:         "",
										Optional:     false,
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
							Optional:     false,
							Skip:         false,
						},
					},

					{
						ItemName: "created_at",
						BaseItem: Scalar{"int", TypeInteger, false},
						Meta: meta.Meta{
							OriginalName: "CreatedAt",
							Name:         "created_at",
							Type:         "",
							Optional:     false,
							Skip:         false,
						},
					},
				},
				Nullable: false,
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
			Expected: List{
				ItemName: "Strings",
				BaseItem: Scalar{"string", TypeString, false},
				Nullable: false,
				Length:   EmptyLength,
			},
		},

		// Ints
		{
			Description: "parse []ints",
			Source:      Ints{},
			Expected: List{
				ItemName: "Ints",
				BaseItem: Scalar{"int", TypeInteger, false},
				Nullable: false,
				Length:   EmptyLength,
			},
		},

		// Floats
		{
			Description: "parse []floats",
			Source:      Floats{},
			Expected: List{
				ItemName: "Floats",
				BaseItem: Scalar{"float32", TypeFloat, false},
				Nullable: false,
				Length:   EmptyLength,
			},
		},

		// Structs
		{
			Description: "parse []structs",
			Source:      Structs{},
			Expected: List{
				ItemName: "Structs",
				BaseItem: Struct{
					ItemName: "CustomType",
					Fields: []Field{
						{
							ItemName: "Name",
							BaseItem: Scalar{"string", TypeString, false},
							Meta: meta.Meta{
								OriginalName: "Name",
								Name:         "Name",
								Type:         "",
								Optional:     false,
								Skip:         false,
							},
						},
					},
				},
				Nullable: false,
				Length:   EmptyLength,
			},
		},

		// String pointers
		{
			Description: "parse []*string",
			Source:      StringPtrs{},
			Expected: List{
				ItemName: "StringPtrs",
				BaseItem: Scalar{"string", TypeString, true},
				Nullable: false,
				Length:   EmptyLength,
			},
		},

		// List of lists
		{
			Description: "parse [][]int",
			Source:      ListList{},
			Expected: List{
				ItemName: "ListList",
				BaseItem: List{
					ItemName: "", // The inner list has no name
					BaseItem: Scalar{"int", TypeInteger, false},
					Length:   EmptyLength,
					Nullable: false,
				},
				Nullable: false,
				Length:   EmptyLength,
			},
		},

		// List pointer
		{
			Description: "parse *[]int",
			Source:      *new(ListPtr), // new(ListPtr) returns a pointer to a nil slice, that is intentionally unhandled by the parser and will return an error for now
			Expected: List{
				ItemName: "",
				BaseItem: Scalar{"int", TypeInteger, false},
				Nullable: true,
				Length:   EmptyLength,
			},
		},

		// List of list pointers
		{
			Description: "parse []*[]int",
			Source:      ListPtrs{},
			Expected: List{
				ItemName: "ListPtrs",
				Length:   EmptyLength,
				BaseItem: List{
					ItemName: "",
					BaseItem: Scalar{"int", TypeInteger, false},
					Length:   EmptyLength,
					Nullable: true,
				},
			},
		},

		// Fixed strings
		{
			Description: "parse [3]string",
			Source:      FixedStrings{},
			Expected: List{
				ItemName: "FixedStrings",
				BaseItem: Scalar{"string", TypeString, false},
				Length:   3,
				Nullable: false,
			},
		},

		// Fixed Structs
		{
			Description: "parse [8]structs",
			Source:      FixedStructs{},
			Expected: List{
				ItemName: "FixedStructs",
				BaseItem: Struct{
					ItemName: "CustomType",
					Fields: []Field{
						{
							ItemName: "Name",
							BaseItem: Scalar{"string", TypeString, false},
							Meta: meta.Meta{
								OriginalName: "Name",
								Name:         "Name",
								Type:         "",
								Optional:     false,
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
			Expected: List{
				ItemName: "FixedIntPtrs",
				BaseItem: Scalar{"int", TypeInteger, true},
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
			Expected: Function{
				ItemName: "Func1",
				Params:   []Item{},
				Returns:  []Item{Scalar{"error", TypeString, false}},
				Nullable: false,
			},
		},

		// Add
		{
			Description: "parse func(int, int) int",
			Source:      Add(nil),
			Expected: Function{
				ItemName: "Add",
				Params: []Item{
					Scalar{"int", TypeInteger, false},
					Scalar{"int", TypeInteger, false},
				},
				Returns:  []Item{Scalar{"int", TypeInteger, false}},
				Nullable: false,
			},
		},

		// ReturnMultiple
		{
			Description: "parse func(string, string) (int, error)",
			Source:      ReturnMultiple(nil),
			Expected: Function{
				ItemName: "ReturnMultiple",
				Params: []Item{
					Scalar{"string", TypeString, false},
					Scalar{"string", TypeString, true},
				},
				Returns: []Item{
					Scalar{"int", TypeInteger, false},
					Scalar{"error", TypeString, false},
				},
				Nullable: false,
			},
		},

		// InsertFoo
		{
			Description: "parse func(Foo) error",
			Source:      InsertFoo(nil),
			Expected: Function{
				ItemName: "InsertFoo",
				Params: []Item{
					Struct{
						ItemName: "Foo",
						Fields: []Field{
							{
								ItemName: "Name",
								BaseItem: Scalar{"string", TypeString, false},
								Meta: meta.Meta{
									OriginalName: "Name",
									Name:         "Name",
									Type:         "",
									Optional:     false,
									Skip:         false,
								},
							},
						},
						Nullable: false,
					},
				},
				Returns:  []Item{Scalar{"error", TypeString, false}},
				Nullable: false,
			},
		},
	}

	runTests(t, tests)
}

func Test_ParseEmbeddedStruct(t *testing.T) {
	type FooEmbedded struct {
		Name string
	}

	type FooParent struct {
		FooEmbedded
		Age int
	}

	tests := []Test{
		{
			Description: "parse embedded struct",
			Source:      FooParent{},
			Expected: Struct{
				ItemName: "FooParent",
				Fields: []Field{
					{
						ItemName: "Name",
						BaseItem: Scalar{"string", TypeString, false},
						Meta: meta.Meta{
							OriginalName: "Name",
							Name:         "Name",
							Type:         "",
							Optional:     false,
							Skip:         false,
						},
					},
					{
						ItemName: "Age",
						BaseItem: Scalar{"int", TypeInteger, false},
						Meta: meta.Meta{
							OriginalName: "Age",
							Name:         "Age",
							Type:         "",
							Optional:     false,
							Skip:         false,
						},
					},
				},
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []Test) {
	for _, tt := range tests {
		runTest(t, tt)
	}
}

func runTest(t *testing.T, tt Test) {
	p := New()

	got, err := p.Parse(reflect.TypeOf(tt.Source))
	if err != nil && !tt.WantErr {
		t.Errorf("[%s] wanted NO error, got error `%s`", tt.Description, err.Error())
	}

	if err == nil && tt.WantErr {
		t.Errorf("[%s] wanted error, got no error", tt.Description)
	}

	if !reflect.DeepEqual(got, tt.Expected) {
		t.Errorf("[%s] wanted %v, got %v", tt.Description, tt.Expected, got)
	}
}
