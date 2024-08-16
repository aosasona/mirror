package typescript

import (
	"errors"
	"fmt"
	"strings"

	"go.trulyao.dev/mirror/config"
	"go.trulyao.dev/mirror/generator"
	"go.trulyao.dev/mirror/parser"
	"go.trulyao.dev/mirror/types"
)

var fileHeader = `/**
* This file was generated by mirror, do not edit it manually as it will be overwritten.
*
* You can find the docs and source code for mirror here: https://github.com/aosasona/mirror
*/
`

type Generator struct {
	config *Config
	indent string
	parser types.ParserInterface
}

func NewGenerator(c *Config) *Generator {
	g := Generator{config: c}

	if c.IndentationType == config.Space {
		g.indent = strings.Repeat(" ", c.IndentationCount)
	} else {
		// 4 spaces to a tab
		g.indent = strings.Repeat("\t", c.IndentationCount/4)
	}

	return &g
}

// SetHeaderText sets the header text for the generated file
func (g *Generator) SetHeaderText(header string) {
	if header == "" {
		return
	}

	fileHeader = header
}

// SetParser sets the parser to use for generating the "types tree"
func (g *Generator) SetParser(parser types.ParserInterface) error {
	if parser == nil {
		return generator.ErrNoParser
	}

	g.parser = parser
	return nil
}

// GenerateItem generates a single item passed to it
func (g *Generator) GenerateItem(item parser.Item) (string, error) {
	var (
		typeString = "type %s = %s"
		baseType   string
		err        error
	)

	baseType, err = g.generateBaseType(item)
	if err != nil {
		return "", err
	}

	if g.config.InludeSemiColon {
		typeString = strings.TrimSpace(typeString) + ";"
	}

	return fmt.Sprintf(typeString, item.Name(), baseType), nil
}

// generateBaseType generates the base type for the item without any additional information
// For example, a scalar type will return `string` or `number` while a list will return `string[]` or `Array<string>`, this is then used by `GenerateItem` to generate the full type definition
func (g *Generator) generateBaseType(item parser.Item, nestingLevel ...int) (string, error) {
	var (
		baseType string
		err      error
	)

	switch item := item.(type) {
	case parser.Scalar:
		baseType, err = g.generateScalar(item)
	case parser.List:
		baseType, err = g.generateList(item)
	case parser.Struct:
		level := 1
		if len(nestingLevel) > 0 {
			level = nestingLevel[0]
		}
		baseType, err = g.generateStruct(item, level)
	case parser.Map:
		baseType, err = g.generateMap(item)
	case parser.Function:
		baseType, err = g.generateFunction(item)
	default:
		return "", generator.ErrUnknownType
	}

	if err != nil {
		return "", err
	}

	if baseType == "" {
		return "", generator.ErrUnknownType
	}

	if item.IsNullable() {
		if g.config.PreferNullForNullable {
			baseType += " | null"
		} else {
			baseType += " | undefined"
		}
	}

	return baseType, nil
}

// GenerateAll implements generator.GeneratorInterface.
func (g *Generator) GenerateAll() ([]string, error) {
	panic("unimplemented")
}

// GenerateN implements generator.GeneratorInterface.
func (g *Generator) GenerateN(int) (string, error) {
	panic("unimplemented")
}

// getScalarRepresentation returns the typescript representation of a scalar type
func (g *Generator) getScalarRepresentation(mirrorType parser.Type) string {
	var typeValue string

	switch mirrorType {
	case parser.TypeAny:
		typeValue = "any"
		if g.config.PreferUnknown {
			typeValue = "unknown"
		}
	case parser.TypeInteger, parser.TypeFloat:
		typeValue = "number"
	case parser.TypeString:
		typeValue = "string"
	case parser.TypeBoolean:
		typeValue = "boolean"
	case parser.TypeByte:
		typeValue = "string"
	// typeValue = "Uint8Array" // TODO: it should most definitely be Uint8Array but I'm not sure if it's a good idea currently with serialization and whatnot
	case parser.TypeTimestamp:
		typeValue = "string"
		// typeValue = "Date" // TODO: generate code to automatically handle this on the TS side if need be in the future like TypeByte
	default:
		return ""
	}

	return typeValue
}

func (g *Generator) generateScalar(item parser.Scalar) (string, error) {
	typeValue := g.getScalarRepresentation(item.Type())
	if typeValue == "" {
		return "", generator.ErrUnknownType
	}

	return typeValue, nil
}

func (g *Generator) generateStruct(item parser.Struct, nestingLevel int) (string, error) {
	var fields []string

	for _, field := range item.Fields {
		if field.Meta.Skip {
			continue
		}

		var (
			fieldName       = field.ItemName
			fieldStr        string
			hasOptionalChar bool
		)

		for i := 0; i < nestingLevel; i++ {
			fieldStr += g.indent
		}

		if field.ItemName == "" && field.Meta.Name == "" {
			return "", generator.ErrNoName
		}

		if field.Meta.Name != "" {
			fieldName = field.Meta.Name
		}

		fieldStr += fieldName

		if field.Meta.Optional {
			fieldStr += "?"
			hasOptionalChar = true
		}

		fieldStr += ": "

		if field.Meta.Type != "" {
			fieldStr += field.Meta.Type
		} else {
			if !g.config.InlineObjects && field.BaseItem.Type() == parser.TypeStruct {
				fieldStr += field.BaseItem.Name()
			} else {
				generatedType, err := g.generateBaseType(field.BaseItem, nestingLevel+1)
				if err != nil {
					return "", err
				}
				fieldStr += generatedType
			}
		}

		// Make sure we don't end up with something like `name?: string | undefined;` as they are equivalent in TS
		if hasOptionalChar && !g.config.PreferNullForNullable {
			fieldStr = strings.TrimSuffix(fieldStr, " | undefined")
		}

		fieldStr += ";"

		fields = append(fields, fieldStr)
	}

	if len(fields) == 0 {
		return "", generator.ErrNoFields
	}

	typeString := "{\n%s\n" + strings.Repeat(g.indent, nestingLevel-1) + "}"
	return fmt.Sprintf(typeString, strings.Join(fields, "\n")), nil
}

func (g *Generator) generateList(item parser.List) (string, error) {
	var (
		listString = ""
		err        error
	)

	if g.config.PreferArrayGeneric {
		listString = "Array<%s>"
	} else {
		listString = "%s[]"
	}

	if item.BaseItem == nil {
		return "", generator.ErrNoBaseItem
	}

	var baseType string
	if item.BaseItem.IsScalar() {
		if baseType, err = g.generateScalar(item.BaseItem.(parser.Scalar)); err != nil {
			return "", err
		}

		if item.BaseItem.IsNullable() {
			if g.config.PreferNullForNullable {
				baseType = fmt.Sprintf("%s | null", baseType)
			} else {
				baseType = fmt.Sprintf("%s | undefined", baseType)
			}

			if !g.config.PreferArrayGeneric {
				baseType = fmt.Sprintf("(%s)", baseType)
			}
		}
	} else {
		baseType = item.BaseItem.Name()
		if g.config.InlineObjects {
			if baseType, err = g.generateBaseType(item.BaseItem); err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf(listString, baseType), nil
}

func (g *Generator) generateMap(item parser.Map) (string, error) {
	typeString := "Record<%s, %s>"

	if item.Key == nil || item.Value == nil {
		return "", generator.ErrNoBaseItem
	}

	var (
		keyType, valueType string
		err                error
	)

	if !item.Key.IsScalar() {
		return "", fmt.Errorf("non-scalar map key (%s) is not supported", item.Key.Name())
	}

	if keyType, err = g.generateScalar(item.Key.(parser.Scalar)); err != nil {
		return "", err
	}

	if valueType, err = g.generateBaseType(item.Value); err != nil {
		return "", err
	}

	return fmt.Sprintf(typeString, keyType, valueType), nil
}

func (g *Generator) generateFunction(item parser.Function) (string, error) {
	var (
		parameterTypes []string
		returnType     = "void"
		err            error
	)

	var paramStr string
	for idx, param := range item.Params {
		if paramStr, err = g.generateBaseType(param); err != nil {
			return "", err
		}

		parameterTypes = append(parameterTypes, "arg"+fmt.Sprint(idx)+": "+paramStr)
	}

	if len(item.Returns) > 1 {
		return "", errors.New("multiple return values are not supported in typescript")
	}

	if len(item.Returns) > 0 {
		returnItem := item.Returns[0]
		if returnType, err = g.generateBaseType(returnItem); err != nil {
			return "", err
		}

		// Surround the return type with parentheses if it's a function
		if returnItem.Type() == parser.TypeFunction {
			returnType = "(" + returnType + ")"
		}
	}

	return fmt.Sprintf("(%s) => %s", strings.Join(parameterTypes, ", "), returnType), nil
}

var _ types.GeneratorInterface = &Generator{}
