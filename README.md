# Mirror 🪞

[![Go Reference](https://pkg.go.dev/badge/go.trulyao.dev/mirror/v2.svg)](https://pkg.go.dev/go.trulyao.dev/mirror/v2)

> [!WARNING]
> This documentation is a work-in-progress and will be updated as required, if you are unclear if mirror supports your usecase or how to do something, please have a look at the GoDocs or open an issue.

Mirror allows you to generate types for other languages using your existing Go types, this can be useful for some form of end-to-end type safety.

# Installation

Mirror can be installed directly using the command below:

```sh
go get go.trulyao.dev/mirror/v2@latest
```

## Usage

The library exposes a large surface area to enable users hook into various parts like the parser and add support for other unsupported languages in user-land, but a normal usage generating typescript types using mirror (with mostly default options) would look like this:

```go
package main

import (
	"fmt"
	"time"

	"go.trulyao.dev/mirror/v2"
	"go.trulyao.dev/mirror/v2/config"
	"go.trulyao.dev/mirror/v2/generator/typescript"
	"go.trulyao.dev/mirror/v2/parser"
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
		Enabled:                os.Getenv("ENV") == "dev", // only enable mirror in dev
		FlattenEmbeddedTypes: false,
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

## Output

```typescript
/**
 * This file was generated by mirror, do not edit it manually as it will be overwritten.
 *
 * You can find the docs and source code for mirror here: https://github.com/aosasona/mirror
 */

type Language = string;

type Address = {
	line_1: string | null;
	line_2: string | null;
	street: string;
	city: string;
	state: string;
	postal_code: string;
	country: string;
};

type Tags = Record<string, string>;

type Person = {
	first_name: string;
	last_name: string;
	age: number;
	address: Address;
	languages: Array<string>;
	grades?: Record<string, number>;
	tags: Record<string, string>;
	created_at: string;
	updated_at: number | null;
	deleted_at: string | null;
	is_active: boolean;
};

type CreateUserFunc = (arg0: Person) => string;
```

See [examples](https://github.com/aosasona/mirror/tree/master/examples) for more options and examples.

## Supported languages

- Typescript
  > More will be added to the library in the future as required

## Tags

You can configure the generated types using struct tags; the `json` tag, the `mirror` tag or the legacy `ts` struct tag. You can pass in the following override values via struct field tags:

- name (string)
- type (string)
- optional (only `true` or `1` or it is ignored)
- skip (only `true` or `1`, but can also simply be written like this: `mirror:"-"`)

#### Example

```go
type Ex struct {
	ID	string `json:"user_id,omitempty" mirror:"type:Uppercase<string>"`
	Name string `mirror:"name:fname,optional:true"`
}
```

This will translate into:

```typescript
export type Ex = {
	user_id?: Uppercase<string> | null;
	fname?: string;
};
```

These give you more control over what types end up being generated. You don't need to specify these, they are optional, if they are not specified, the default values are inferred from the types themselves.

## Contribution

PRs and issues are welcome :)

## Development

- To run the example:

```sh
just example
```

- To run the tests:

```sh
just test
```
