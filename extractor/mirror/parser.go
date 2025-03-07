package mirrormeta

import (
	"fmt"
	"log/slog"
	"strings"

	mt "go.trulyao.dev/mirror/v2/extractor/meta"
)

func ref[T any](t T) *T {
	return &t
}

func deref[T any](t *T, fallback ...T) T {
	var f T

	if len(fallback) > 0 {
		f = fallback[0]
	}

	if t == nil {
		return f
	}

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

const (
	BoolErrorMessage = "Using numbers as boolean values is deprecated, please use 'true' or 'false' instead, this will be removed in the future"
)

type ParsedMeta struct {
	Name     *string
	Type     *string
	Skip     *bool
	Optional mt.Optional
}

type MetaParser struct {
	// The input string to parse (as a slice of runes)
	input string
}

func NewMetaParser(input string) *MetaParser {
	input = strings.TrimSpace(input)
	return &MetaParser{input}
}

// Parse the input string and return a ParsedMeta struct
// This methods parses other fields first and then attempts to use whatever is left as the type as a process of elimination
func (p *MetaParser) Parse() (*ParsedMeta, error) {
	var meta ParsedMeta

	// Find a valid `optional:1|0|true|false` part of the string
	optional, err := p.parseBool("optional")
	if err != nil {
		return nil, err
	}
	if optional != nil {
		if *optional {
			meta.Optional = mt.OptionalTrue
		} else {
			meta.Optional = mt.OptionalFalse
		}
	}

	// Find a valid `skip:1|0|true|false` part of the string
	skip, err := p.parseBool("skip")
	if err != nil {
		return nil, err
	}
	if skip != nil {
		meta.Skip = skip
	}

	// Find a valid `name:...` part of the string, valid names are alphanumeric, underscores and dashes
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if name != nil {
		meta.Name = name
	}

	// The remaining part of the string should ideally be the type
	if strings.Contains(p.input, "type:") {
		t := trimInvalidChars(strings.Replace(p.input, "type:", "", 1))
		if t != "" {
			meta.Type = ref(t)
		}
	}

	return &meta, nil
}

func (p *MetaParser) parseName() (*string, error) {
	// We need to manually parse this to make sure we don't end up with an attribute that is in a block, for example: {optional:1}
	var (
		// This will be used to keep track of whether we are in a block or not
		// For instance, if we hit a '{' n times, we need to make sure we hit a '}' n times before we can say we are out of the block
		blocks []rune

		rangeStart int // this will be used to store the start of the attribute
		value      string

		attrLen = len("name:")
	)

	for i, r := range p.input {
		if isBlockStart(r) {
			blocks = append(blocks, r)
		} else if isBlockEnd(r) {
			// If this does not close the matching block in descending order, it is an error
			// For example, if this was our block stack: ['{', '[', '('] and we hit a '}', it is an error, we need to hit a ')' first
			// If the block stack is empty, it means we have an extra closing block, which is an error as well
			if len(blocks) == 0 || !isMatchingBlock(blocks[len(blocks)-1], r) {
				return nil, fmt.Errorf("unexpected block closing: %c", r)
			}

			// Pop the last block
			blocks = blocks[:len(blocks)-1]
		} else if len(blocks) == 0 && isAlphaNumeric(r) {
			// Check if we have an optional attribute
			if p.readRange(i, attrLen) == "name:" {
				rangeStart = i

				// Read till the end of the name, that is, we have run out of alphanumeric characters (including underscores and dashes)
				for _, r := range p.input[i+attrLen:] {
					if !isAlphaNumeric(r) {
						break
					}

					value += string(r)
				}
			}
		}
	}

	if value == "" {
		return nil, nil
	}

	p.truncateRange(rangeStart, attrLen+len(value))
	return ref(value), nil
}

func (p *MetaParser) parseBool(attribute string) (*bool, error) {
	attribute = attribute + ":"

	// We need to manually parse this to make sure we don't end up with an attribute that is in a block, for example: {optional:1}
	var (
		// This will be used to keep track of whether we are in a block or not
		// For instance, if we hit a '{' n times, we need to make sure we hit a '}' n times before we can say we are out of the block
		blocks []rune

		rangeStart int // this will be used to store the start and end of the optional attribute
		value      string

		attrLen = len(attribute)
	)

	for i, r := range p.input {
		if isBlockStart(r) {
			blocks = append(blocks, r)
			continue
		} else if isBlockEnd(r) {
			// If this does not close the matching block in descending order, it is an error
			// For example, if this was our block stack: ['{', '[', '('] and we hit a '}', it is an error, we need to hit a ')' first
			// If the block stack is empty, it means we have an extra closing block, which is an error as well
			if len(blocks) == 0 || !isMatchingBlock(blocks[len(blocks)-1], r) {
				return nil, fmt.Errorf("unexpected block closing while parsing `%s`: %c", attribute, r)
			}

			// Pop the last block
			blocks = blocks[:len(blocks)-1]
		} else if len(blocks) == 0 && isAlphaNumeric(r) {
			// Check if we have an optional attribute
			if p.readRange(i, attrLen) == attribute {
				rangeStart = i

				// Check if the attribute is set to true or false
				if v := p.readRange(i+attrLen, 4); v == "true" {
					value = v
					break
				} else if v := p.readRange(i+attrLen, 5); v == "false" {
					value = v
					break
				} else if v := p.readRange(i+attrLen, 1); v == "1" {
					slog.Debug(BoolErrorMessage)
					value = v
					break
				} else if v := p.readRange(i+attrLen, 1); v == "0" {
					slog.Debug(BoolErrorMessage)
					value = v
					break
				}
			}
		}
		// We don't care about other characters, hence we ignore them
	}

	// If a block was opened and not closed, it is an error
	if len(blocks) > 0 {
		return nil, fmt.Errorf(
			"unexpected block opening while parsing `%s`: %c",
			attribute,
			blocks[len(blocks)-1],
		)
	}

	if value == "" {
		return nil, nil
	}

	p.truncateRange(rangeStart, attrLen+len(value))
	booleanValue := value == "true" || value == "1"
	return ref(booleanValue), nil
}

func (p *MetaParser) readRange(start int, length int) string {
	if start+length > len(p.input) {
		return ""
	}
	return p.input[start : start+length]
}

// This is to remove sections of the input string that have already been parsed
func (p *MetaParser) truncateRange(start int, length int) {
	p.input = strings.TrimSpace(p.input[:start] + p.input[start+length:])
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' ||
		r == '-'
}

func isBlockStart(r rune) bool {
	return r == '{' || r == '[' || r == '(' || r == '<' || r == '"'
}

func isBlockEnd(r rune) bool {
	return r == '}' || r == ']' || r == ')' || r == '>' || r == '"'
}

func isMatchingBlock(start, end rune) bool {
	return (start == '{' && end == '}') || (start == '[' && end == ']') ||
		(start == '(' && end == ')') ||
		(start == '<' && end == '>')
}

func trimInvalidChars(s string) string {
	var value string
	for i := len(s) - 1; i >= 0; i-- {
		if !isWhitespace(rune(s[i])) && s[i] != ',' {
			// return s[:i+1]
			value = s[:i+1]
			break
		}
	}

	for i := 0; i < len(value); i++ {
		if !isWhitespace(rune(value[i])) && value[i] != ',' {
			return value[i:]
		}
	}

	return ""
}
