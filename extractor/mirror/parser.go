package mirrormeta

import "fmt"

func ref[T any](t T) *T {
	return &t
}

func deref[T any](t *T) T {
	return *t
}

type Attribute int

const (
	AtributeNone Attribute = iota
	AtributeName
	AtributeType
	AttributeSkip
	AtributeOptional
)

type ParsedMeta struct {
	Name     *string
	Type     *string
	Skip     *bool
	Optional *bool
}

type MetaParser struct {
	// The input string to parse (as a slice of runes)
	chars []rune

	// The current attribute being parsed
	currentAttribute Attribute

	// The current index of the input string
	currentIndex int

	// The current attribute being parsed
	node ParsedMeta
}

func NewMetaParser(input string) *MetaParser {
	return &MetaParser{
		chars:            []rune(input),
		currentAttribute: AtributeNone,
		currentIndex:     0,
		node:             ParsedMeta{},
	}
}

func (p *MetaParser) Parse() (*ParsedMeta, error) {
	for p.currentIndex < len(p.chars) {
		switch {
		case isWhitespace(p.char()):
			p.next()

		// Read the attribute name - this put us in a parsing state
		case isAlphaNumeric(p.char()) && p.currentAttribute == AtributeNone:
			identifier := p.readIdentifier()
			p.currentAttribute = p.getAttributeFromString(identifier)

		default:
			// If it is not an alphanumeric character and we are not parsing a node already, then it is an error
			if p.currentAttribute == AtributeNone {
				return nil, fmt.Errorf(
					"unexpected character: %c, expected one of `type`, `name`, `skip`, `optional`",
					p.char(),
				)
			}

			var value string
			// This tracks whether we are in a block like `type:{ ... }`
			for p.currentIndex < len(p.chars) {
				// Handle the type attribute first as it is the most complex bit
				if p.currentAttribute == AtributeType {
					// If it is a { [, (, or <, then we are in a block, read till we find the closing block
					if isBlockStart(p.char()) {
						value = p.readBlock()
						break
					}

					// If it is not a block, then read till we find a comma or whitespace
					start := p.currentIndex
					for !isWhitespace(p.char()) && p.char() != ',' {
						p.next()
					}

					value = string(p.chars[start:p.currentIndex])

					// If we find a comma, then we are done reading the value
					if p.char() == ',' {
						p.next()
						break
					}
				}
			}
		}
	}

	panic("TODO: implement")
}

func (p *MetaParser) readIdentifier() string {
	start := p.currentIndex
	for isAlphaNumeric(p.char()) {
		p.next()
	}

	return string(p.chars[start:p.currentIndex])
}

func (p *MetaParser) readBlock() string {
	start := p.currentIndex

	for p.char() != '}' && p.char() != ']' && p.char() != ')' && p.char() != '>' {
		p.next()
	}

	return string(p.chars[start:p.currentIndex])
}

func (p *MetaParser) char() rune {
	if p.currentIndex >= len(p.chars) {
		return 0
	}

	return p.chars[p.currentIndex]
}

func (p *MetaParser) peek() rune {
	if p.currentIndex+1 >= len(p.chars) {
		return 0
	}

	return p.chars[p.currentIndex+1]
}

func (p *MetaParser) next() rune {
	if p.currentIndex >= len(p.chars) {
		return 0
	}

	p.currentIndex++
	return p.chars[p.currentIndex]
}

func (p *MetaParser) getAttributeFromString(text string) Attribute {
	var attribute Attribute
	switch text {
	case "name":
		attribute = AtributeName
	case "type":
		attribute = AtributeType
	case "skip":
		attribute = AttributeSkip
	case "optional":
		attribute = AtributeOptional
	}

	return attribute
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func isBlockStart(r rune) bool {
	return r == '{' || r == '[' || r == '(' || r == '<'
}
