package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"go.trulyao.dev/mirror/v2"
	"go.trulyao.dev/mirror/v2/config"
	"go.trulyao.dev/mirror/v2/extractor/meta"
	"go.trulyao.dev/mirror/v2/generator/typescript"
	"go.trulyao.dev/mirror/v2/parser"
)

type (
	Language string
	Tags     map[string]string
	NamedAny any
)

type Address struct {
	Line1      *string `json:"line_1"      mirror:"type:string"`
	Line2      *string `json:"line_2"`
	Street     string  `                   mirror:"name:street"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"`
}

type Person struct {
	FName     string `mirror:"name:first_name"`
	LName     string `mirror:"name:last_name"`
	Age       int    `mirror:"name:age"`
	Address   `mirror:"name:address"`
	Languages []Language     `mirror:"name:languages"`
	Grades    map[string]int `mirror:"name:grades,optional:1"`
	Tags      Tags           `mirror:"name:tags"`
	Props     any            `mirror:"name:props,optional:true"`
	CreatedAt time.Time      `mirror:"name:created_at"`
	UpdatedAt *time.Time     `mirror:"name:updated_at,type:number"`
	DeletedAt *time.Time     `mirror:"name:deleted_at"`
	IsActive  bool           `mirror:"name:is_active"`
	Error     error          `mirror:"name:error"`
}

type StateMeta struct {
	ExpiresAt time.Time `mirror:"name:expires_at"`
	CreatedAt time.Time `mirror:"name:created_at"`
	Meta      NamedAny  `mirror:"name:meta"`
	User      NamedAny  `mirror:"type:{ user_id: string, role: 'admin' | 'user', tags: Array<string> }" json:"user"`
}

type Store struct {
	Key   string    `mirror:"name:key"`
	Value string    `mirror:"name:value"`
	Meta  StateMeta `mirror:"name:meta"`
}

type UserWithNestedProperties struct {
	FirstName string `mirror:"name:first_name"`
	LastName  string `mirror:"name:last_name"`

	// These are to test nesting and indentation
	Stores     []Store          `mirror:"name:stores"`
	OtherStore map[string]Store `mirror:"name:other_store"`
}

type Collection struct {
	Items []string `mirror:"name:items"`
	Desc  string   `mirror:"name:desc"`
}

type CreateUserFunc func(p Person) error

func main() {
	start := time.Now()

	m := mirror.New(config.Config{
		Enabled:              true,
		FlattenEmbeddedTypes: false,
	})

	m.Parser().
		OnParseField(func(_ *reflect.Type, _ *reflect.StructField, field *parser.Field) error {
			// Rename the `desc` field to `description` and make it optional
			if field.ItemName == "desc" {
				field.Meta.Name = "description"
				field.ItemName = "description"
				field.Meta.Optional = meta.OptionalTrue
			}

			return nil
		})

	m.Parser().OnParseItem(func(sourceName string, target parser.Item) error {
		// Add a new `created_at` field to the `Collection` struct
		if sourceName == "Collection" {
			createdAtField := parser.Field{
				ItemName: "CreatedAt",
				BaseItem: &parser.Scalar{
					ItemName: "",
					ItemType: parser.TypeString,
				},
				Meta: meta.Meta{Name: "created_at", Type: "Date"},
			}

			if target, ok := target.(*parser.Struct); ok {
				target.Fields = append(target.Fields, createdAtField)
			}
		}

		return nil
	})

	m.AddSources(
		Language(""),
		Address{},
		Tags{},
		Person{},
		Store{},
		StateMeta{},
		UserWithNestedProperties{},
		Collection{},
		CreateUserFunc(nil),
	)

	defaultTS := typescript.DefaultConfig().
		SetFileName("default.ts").
		SetOutputPath("./examples").
		SetIndentationType(config.IndentTab)

	inlinedTS := typescript.DefaultConfig().
		SetFileName("inlined.ts").
		SetOutputPath("./examples").
		SetInlineObjects(true).
		SetPrefix("Inline_").
		SetIndentationType(config.IndentTab)

	m.AddTarget(defaultTS).AddTarget(inlinedTS)

	err := m.GenerateAndSaveAll()
	if err != nil {
		fmt.Println(err)
	}

	// Flatten embedded structs
	newParser := parser.New()
	flattenedTS := typescript.DefaultConfig().
		SetFileName("flattened"). // no extension
		SetOutputPath("./examples")

	m.ResetTargets().
		SetParser(newParser).
		AddTarget(flattenedTS).
		AddSources(*new(Language), Tags{}, Person{}, Collection{}, CreateUserFunc(nil))

	if err = newParser.SetConfig(parser.Config{
		EnableCaching:        true,
		FlattenEmbeddedTypes: true,
	}); err != nil {
		log.Fatal(err)
	}

	err = m.GenerateAndSaveAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated %d types in %s\n", m.Count(), time.Since(start))
}
