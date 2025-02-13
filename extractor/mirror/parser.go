package mirrormeta

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
	panic("TODO: implement")
}

func (p *MetaParser) setCurrentAttribute(text string) {
	switch text {
	case "name":
		p.currentAttribute = AtributeName
	case "type":
		p.currentAttribute = AtributeType
	case "skip":
		p.currentAttribute = AttributeSkip
	case "optional":
		p.currentAttribute = AtributeOptional
	}
}
