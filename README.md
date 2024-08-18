# Mirror 🪞

[![Go Reference](https://pkg.go.dev/badge/go.trulyao.dev/mirror.svg)](https://pkg.go.dev/go.trulyao.dev/mirror)

Mirror allows you to generate types for other languages using your existing Go types. V1 only had support for Golang-to-Typescript types

> [!NOTE]
> This package was previously known as [`gots`](https://github.com/aosasona/gots).

# Installation

Mirror can be installed directly using the command below:

```sh
go get go.trulyao.dev/mirror@latest
```

## Usage

Mirror supports generating types for other languages but only has limited support built-in right now for very few languages. The library exposes a large surface area to enable users hook into various parts like the parser and add support for other unsupported languages in user-land. Generating typescript types using mirror (+ default options) would look like this:

```go
package main

import (
	"fmt"
	"time"

	"go.trulyao.dev/mirror"
	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/generator/typescript"
	"go.trulyao.dev/mirror/parser"
)

type Language string

type Tags map[string]string

type Address struct {
	Line1      *string `mirror:"name:line_1"`
	Line2      *string `mirror:"name:line_2"`
	Street     string  `mirror:"name:street"`
	City       string  `mirror:"name:city"`
	State      string  `mirror:"name:state"`
	PostalCode string  `mirror:"name:postal_code"`
	Country    string  `mirror:"name:country"`
}

type Person struct {
	FName     string `mirror:"name:first_name"`
	LName     string `mirror:"name:last_name"`
	Age       int    `mirror:"name:age"`
	Address   `mirror:"name:address"`
	Languages []Language     `mirror:"name:languages"`
	Grades    map[string]int `mirror:"name:grades,optional:1"`
	Tags      Tags           `mirror:"name:tags"`
	CreatedAt time.Time      `mirror:"name:created_at"`
	UpdatedAt *time.Time     `mirror:"name:updated_at,type:number"`
	DeletedAt *time.Time     `mirror:"name:deleted_at"`
	IsActive  bool           `ts:"name:is_active"` // using deprecated `ts` tag
}

type CreateUserFunc func(p Person) error

func main() {
	m := mirror.New(config.Config{
		Enabled:                true,
		FlattenEmbeddedStructs: false,
	})

	m.AddSources(
	  Language(""),
	  Address{},
	  Tags{},
	  Person{},
	  CreateUserFunc(nil),
	)

	target := typescript.DefaultConfig().
		SetOutputPath(".").
		SetFileName("generated.ts").
		SetIndentationType(config.IndentTab)

	m.AddTarget(target)

	if err := m.GenerateAndSaveAll(); err != nil {
		fmt.Println(err)
	}
}
```

# Supported languages

- Typescript
  > More will be added to the library in the future as required
