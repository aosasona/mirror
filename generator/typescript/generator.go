package typescript

import (
	"errors"
	"fmt"
	"strings"

	"go.trulyao.dev/mirror/v2/config"
	"go.trulyao.dev/mirror/v2/extractor/meta"
	"go.trulyao.dev/mirror/v2/parser"
	"go.trulyao.dev/mirror/v2/types"
)

var fileHeader = `/**
 * This file was generated by mirror, do not edit it manually as it will be overwritten.
 *
 * You can find the docs and source code for mirror here: https://github.com/aosasona/mirror
 */
`

type Generator struct {
	// config is the configuration for the generator
	config *Config

	// indent is the indentation string used internally by the generator
	indent string

	// parser is the parser used to generate the types
	parser types.ParserInterface

	// nonStrict is a flag to determine if the generator should be non-strict
	nonStrict bool
}

// NewGenerator returns a new typescript generator instance with the provided config
func NewGenerator(c *Config) *Generator {
	g := Generator{config: c}

	if c.IndentationType == config.IndentSpace {
		g.indent = strings.Repeat(" ", c.IndentationCount)
	} else {
		// 4 spaces to a tab
		g.indent = strings.Repeat("\t", c.IndentationCount/4)
	}

	return &g
}

// SetNonStrict sets the generator to be non-strict, meaning it will not throw an error if a referenced type does not exist and other strict checks
func (g *Generator) SetNonStrict(strict bool) {
	g.nonStrict = strict
}

// SetHeaderText sets the header text for the generated file
func (g *Generator) SetHeaderText(header string) {
	fileHeader = header
}

// SetParser sets the parser to use for generating the "types tree"
func (g *Generator) SetParser(parser types.ParserInterface) error {
	if parser == nil {
		return errors.New("parser cannot be nil")
	}

	g.parser = parser
	return nil
}

// GenerateItem generates a single item passed to it
func (g *Generator) GenerateItem(item parser.Item) (string, error) {
	var (
		typeString = "export type %s = %s"
		baseType   string
		err        error
	)

	baseType, err = g.generateBaseType(item, nil)
	if err != nil {
		return "", err
	}

	if g.config.InludeSemiColon {
		typeString = strings.TrimSpace(typeString) + ";"
	}

	typeName := item.Name()
	if g.config.TypePrefix != "" {
		typeName = g.config.TypePrefix + typeName
	}

	return fmt.Sprintf(typeString, typeName, baseType), nil
}

// GenerateItemType generate ONLY the type definition for an item (e.g. "string", "{ foo: Bar, ...}")
func (g *Generator) GenerateItemType(item parser.Item) (string, error) {
	var (
		itemType string
		err      error
	)

	if itemType, err = g.generateBaseType(item, nil); err != nil {
		return "", err
	}

	return itemType, nil
}

// generateBaseType generates the base type for the item without any additional information
// For example, a scalar type will return `string` or `number` while a list will return `string[]` or `Array<string>`, this is then used by `GenerateItem` to generate the full type definition
func (g *Generator) generateBaseType(
	item parser.Item,
	metadata *meta.Meta,
	nestingLevel ...int,
) (string, error) {
	var (
		baseType string
		err      error
	)

	level := 1
	if len(nestingLevel) > 0 {
		level = nestingLevel[0]
	}

	switch item := item.(type) {
	case *parser.Scalar:
		baseType, err = g.generateScalar(item)
	case *parser.List:
		baseType, err = g.generateList(item, level)
	case *parser.Struct:
		baseType, err = g.generateStruct(item, level)
	case *parser.Map:
		baseType, err = g.generateMap(item, level)
	case *parser.Function:
		baseType, err = g.generateFunction(item)
	default:
		return "", fmt.Errorf("unknown type: %T", item)
	}

	if err != nil {
		return "", err
	}

	if baseType == "" {
		return "", errors.New("failed to generate base type")
	}

	// Check meta for nullability and optionality
	var optional meta.Optional
	if metadata != nil {
		optional = metadata.Optional
	}

	isOptional := item.IsNullable() && optional.None()
	isOverrideOptional := optional.True()
	if isOptional || isOverrideOptional {
		if g.config.PreferNullForNullable {
			baseType += " | null"
		} else {
			baseType += " | undefined"
		}
	}

	return baseType, nil
}

// GenerateAll generates all the type definitions in the parser
// This method uses the parser's Iterate method to iterate over all the items in the parser without consuming them
func (g *Generator) GenerateAll() ([]string, error) {
	var types []string

	generateTS := func(item parser.Item) error {
		typeDef, err := g.GenerateItem(item)
		if err != nil {
			return err
		}

		types = append(types, typeDef)
		return nil
	}

	if err := g.parser.Iterate(generateTS); err != nil {
		return nil, err
	}

	return types, nil
}

// GenerateN generates the type definition for the nth item in the parser, this operation is 0-indexed and cached by default (unless disabled in the parser)
func (g *Generator) GenerateN(idx int) (string, error) {
	source, err := g.parser.ParseN(idx)
	if err != nil {
		return "", err
	}

	return g.GenerateItem(source)
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
	case parser.TypeTimestamp:
		typeValue = "string"

	// No-oop types
	case parser.TypeVoid:
		typeValue = "void"
	case parser.TypeNil:
		typeValue = "null"

	default:
		return ""
	}

	return typeValue
}

// generateScalar generates the typescript representation of a scalar type (string, number, boolean, etc)
func (g *Generator) generateScalar(item *parser.Scalar) (string, error) {
	typeValue := g.getScalarRepresentation(item.Type())
	if typeValue == "" {
		return "", fmt.Errorf("unknown scalar type: %s", item.Name())
	}

	return typeValue, nil
}

// generateStruct generates the typescript representation of a struct
func (g *Generator) generateStruct(item *parser.Struct, nestingLevel int) (string, error) {
	var fields []string

	for _, field := range item.Fields {
		// Skip fields that are marked to be skipped so they don't appear in the generated types
		if field.Meta.Skip {
			continue
		}

		var (
			fieldName       = field.ItemName
			fieldStr        string
			hasOptionalChar bool
		)

		// Properly indent the fields
		for i := 0; i < nestingLevel; i++ {
			fieldStr += g.indent
		}

		// If the field has no name, we can't generate a type for it
		if field.ItemName == "" && field.Meta.Name == "" {
			return "", fmt.Errorf(
				"unable to find name for field `%s` in struct `%s`",
				field.BaseItem.Name(),
				item.Name(),
			)
		}

		if field.Meta.Name != "" {
			fieldName = field.Meta.Name
		}

		fieldStr += fieldName

		if field.Meta.Optional.True() {
			fieldStr += "?"
			hasOptionalChar = true
		}

		fieldStr += ": "

		// if the field has an override type (using the `mirror` tag), use that
		if field.Meta.Type != "" {
			fieldStr += field.Meta.Type
			isOptional := field.BaseItem.IsNullable() && field.Meta.Optional.None()
			isOverrideOptional := field.Meta.Optional.True()

			if isOptional || isOverrideOptional {
				if g.config.PreferNullForNullable {
					fieldStr += " | null"
				} else {
					fieldStr += " | undefined"
				}
			}
		} else {
			if !g.config.InlineObjects && field.BaseItem.Type() == parser.TypeStruct {
				// Ensure the referenced type exists before proceeding - this is only necessary if inline objects are disabled since we don't want to reference a type that doesn't exist
				if !g.referenceExists(field.BaseItem.Name()) {
					return "", fmt.Errorf("referenced type `%s` does not exist, you need to either enable inline objects or pass in the referenced type", field.BaseItem.Name())
				}

				fieldStr += field.BaseItem.Name()
			} else {
				// Generate the base type for the field
				generatedType, err := g.generateBaseType(field.BaseItem, &field.Meta, nestingLevel+1)
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

	typeString := "{\n%s\n" + strings.Repeat(g.indent, nestingLevel-1) + "}"
	return fmt.Sprintf(typeString, strings.Join(fields, "\n")), nil
}

// generateList generates the typescript representation of a list type (array or slice in Go)
func (g *Generator) generateList(item *parser.List, nestingLevel int) (string, error) {
	var (
		listString string
		err        error
	)

	if g.config.PreferArrayGeneric {
		listString = "Array<%s>"
	} else {
		listString = "%s[]"
	}

	if item.BaseItem == nil {
		return "", fmt.Errorf("no base item found for list type: `%s`", item.Name())
	}

	var baseType string

	// Scalar types are expanded to their types (e.g. string, number, etc) by default
	if item.BaseItem.IsScalar() {
		if baseType, err = g.generateScalar(item.BaseItem.(*parser.Scalar)); err != nil {
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
		// Ensure the referenced type exists before proceeding
		if !g.config.InlineObjects && !g.referenceExists(item.BaseItem.Name()) {
			return "", fmt.Errorf("referenced type `%s` does not exist, you need to either enable inline objects or pass in the referenced type", item.BaseItem.Name())
		}

		// If inline objects are enabled, generate the base type for the item
		baseType = item.BaseItem.Name()
		if g.config.InlineObjects {
			if baseType, err = g.generateBaseType(item.BaseItem, nil, nestingLevel); err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf(listString, baseType), nil
}

// generateMap generates the typescript representation of a map
func (g *Generator) generateMap(item *parser.Map, nestingLevel int) (string, error) {
	typeString := "Record<%s, %s>"

	if item.Key == nil || item.Value == nil {
		return "", fmt.Errorf("key or value is nil for map type: `%s`", item.Name())
	}

	var (
		keyType, valueType string
		err                error
	)

	if !item.Key.IsScalar() {
		return "", fmt.Errorf("non-scalar map key (%s) is not supported", item.Key.Name())
	}

	key, ok := item.Key.(*parser.Scalar)
	if !ok {
		return "", fmt.Errorf("non-scalar map key (%s) is not supported", item.Key.Name())
	}

	if keyType, err = g.generateScalar(key); err != nil {
		return "", err
	}

	if valueType, err = g.generateBaseType(item.Value, nil, nestingLevel); err != nil {
		return "", err
	}

	return fmt.Sprintf(typeString, keyType, valueType), nil
}

// generateFunction generates the typescript representation of a function
func (g *Generator) generateFunction(item *parser.Function) (string, error) {
	var (
		parameterTypes []string
		returnType     = "void"
		err            error
	)

	for idx, param := range item.Params {
		var paramStr string

		// Scalar types are always expanded to their types (e.g. string, number, etc) by default
		if g.config.InlineObjects || param.IsScalar() {
			if paramStr, err = g.generateBaseType(param, nil); err != nil {
				return "", err
			}
		} else {
			paramStr = param.Name()
		}

		parameterTypes = append(parameterTypes, "arg"+fmt.Sprint(idx)+": "+paramStr)
	}

	if len(item.Returns) > 1 {
		return "", errors.New("multiple return values are not supported in typescript")
	}

	if len(item.Returns) > 0 {
		returnItem := item.Returns[0]
		if returnType, err = g.generateBaseType(returnItem, nil); err != nil {
			return "", err
		}

		// Surround the return type with parentheses if it's a function
		if returnItem.Type() == parser.TypeFunction {
			returnType = "(" + returnType + ")"
		}
	}

	return fmt.Sprintf("(%s) => %s", strings.Join(parameterTypes, ", "), returnType), nil
}

// referenceExists() checks if the type being referenced exists in the parser, especially for non-inlined objects
func (g *Generator) referenceExists(name string) bool {
	if g.nonStrict {
		return true
	}

	_, exists := g.parser.LookupByName(name)
	return exists
}
