# Mirror 🪞

> [!NOTE]
> This package was previously known as [`gots`](https://github.com/aosasona/gots). Same code, different name, and development will continue here.

[![Go Reference](https://pkg.go.dev/badge/go.trulyao.dev/mirror.svg)](https://pkg.go.dev/go.trulyao.dev/mirror)

Generate Typescript types from Go types during runtime.

View generated example [here](./examples/example.ts)

## You should know...

The generated types may not always match what you expect (especially embedded structs) and might just be an `any` or `unknown`, to be more specific, it is advised to use the type property in the `mirror` or `ts` struct tag to specify the type yourself. Mirror is not designed or built to be or ever be 100% accurate, just enough to have you setup and ready to communicate with your Go service/app/API _safely_ in Typescript, knowing a large part of what to send and expect back.

# Installation

Just paste this in your terminal (I promise it's safe):

```bash
go get -u go.trulyao.dev/mirror
```

# Usage

Not an exceptional documentation but this should help you get started

```go
package main

import (
	"fmt"
	"time"

	"go.trulyao.dev/mirror"
	"go.trulyao.dev/mirror/config"
)

type Language string

type Tags map[string]string

type Person struct {
	FName     string         `mirror:"name:first_name"`
	LName     string         `mirror:"name:last_name"`
	Age       int            `mirror:"name:age"`
	Languages []Language     `mirror:"name:languages"`
	Grades    map[string]int `mirror:"name:grades,optional:1"`
	Tags      Tags           `mirror:"name:tags"`
	CreatedAt time.Time      `mirror:"name:created_at"`
	UpdatedAt *time.Time     `mirror:"name:updated_at"`
	DeletedAt *time.Time     `mirror:"name:deleted_at"`
}

func main() {
	gt := mirror.New(config.Config{
		Enabled:           mirror.Bool(true),
		OutputFile:        mirror.String("./examples/example.ts"),
		UseTypeForObjects: mirror.Bool(true),
		ExpandObjectTypes: mirror.Bool(true),
	})

	// ===> Individually
	gt.AddSource(*new(Language))
	gt.AddSource(Tags{})
	gt.AddSource(Person{})

	out, err := gt.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}

	// save to file
	err = gt.Commit(out)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ===> As a group
	gt.Register(*new(Language), Tags{}, Person{})

	// generate types and save to the file
	err := gt.Execute(true)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

It is safer to enable mirror in development only, you can do this however way you want in your application. For example:

```go
...
ts := mirror.New(mirror.Config{
	Enabled: mirror.Bool(os.Getenv("ENV") == "development"),
})
...
```

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
export interface Ex {
	user_id?: Uppercase<string>;
	fname?: string;
}
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
