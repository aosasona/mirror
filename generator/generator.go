package generator

import (
	"fmt"
	"reflect"
	"strings"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/parser"
)

// acts as a tab "character" for formatting purposes in the generated code (4 spaces)
const TAB = "    "

var tg *TypeGenerator

type CodeGenerator struct {
	allowUnexportedFields bool
	useTypeForObjects     bool
}

type Opts struct {
	AllowUnexportedFields bool // allow or disallow unexported fields
	UseTypeForObjects     bool // use type instead of interface for objects
	ExpandStructs         bool // whether to expand struct types into object types like { foo: string } instead of the name
	PreferUnknown         bool // whether to prefer unknown over any
}

func NewGenerator(opts Opts) *CodeGenerator {
	tg = NewTypeGenerator(TypeGeneratorOpts{
		ExpandStruct:  opts.ExpandStructs,
		PreferUnknown: opts.PreferUnknown,
	})

	return &CodeGenerator{
		allowUnexportedFields: opts.AllowUnexportedFields,
		useTypeForObjects:     opts.UseTypeForObjects,
	}
}

func (g *CodeGenerator) Generate(src any) string {
	var (
		srcType = reflect.TypeOf(src)
		result  string
	)

	if srcType.Kind() == reflect.Struct {
		result = "export interface %s "
		if g.useTypeForObjects {
			result = "export type %s = "
		}

		result += g.generateObjectType(srcType)
	} else {
		result = "export type %s = "
		result += string(tg.GetFieldType(reflect.StructField{
			Name: srcType.Name(),
			Type: srcType,
		}))
	}

	return fmt.Sprintf(result, srcType.Name())
}

func (g *CodeGenerator) generateObjectType(src reflect.Type) string {
	var result string

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)

		tag, err := parser.Parse(field)
		if err != nil {
			if config.Debug {
				fmt.Printf("Error parsing field: %s\n", err.Error())
			}
			continue
		}

		if (!field.IsExported() && !g.allowUnexportedFields) || tag.Skip {
			continue
		}

		p := makeProperty(field, tg, tag)

		result += fmt.Sprintf("%s%s%s: %s;\n", TAB, p.Name, p.OptionalChar, p.Type)
	}

	return fmt.Sprintf("{\n%s%s\n}", TAB, strings.TrimSpace(result))
}
