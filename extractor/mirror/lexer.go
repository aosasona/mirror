package mirrormeta

type (
	tokenType int
	token     struct {
		// The type of token
		tokenType tokenType

		// The value of the token (if it is a string)
		value string
	}
)

func newToken(tokenType tokenType, value string) token {
	return token{tokenType: tokenType, value: value}
}

func (t token) Value() string { return t.value }

const (
	// tokenTypeLiteral represents any other strings, we do not care about angle brackets, braces, etc
	tokenTypeLiteral tokenType = iota
	tokenTypeComma
	tokenTypeColon
	tokenTypeEOF
)

type lexer struct {
	chars []rune
	index int
}

func newLexer(input string) *lexer {
	return &lexer{
		chars: []rune(input),
		index: 0,
	}
}

func (l *lexer) lex() token {
	for l.index < len(l.chars) {
		l.consumeWhitespace()

		switch {
		case l.next() == ',':
			l.index++
			return newToken(tokenTypeComma, ",")
		case l.next() == ':':
			l.index++
			return newToken(tokenTypeColon, ":")
		case l.next() == 0:
			return newToken(tokenTypeEOF, "")
		default:
			return l.readLiteral()
		}
	}

	return newToken(tokenTypeEOF, "")
}

func (l *lexer) next() rune {
	if l.index >= len(l.chars) {
		return 0
	}

	return l.chars[l.index]
}

func (l *lexer) consumeWhitespace() {
	for l.next() == ' ' || l.next() == '\t' || l.next() == '\n' {
		l.index++
	}
}

func (l *lexer) readLiteral() token {
	start := l.index
	for l.next() != ' ' && l.next() != '\t' && l.next() != '\n' && l.next() != ',' && l.next() != ':' && l.next() != 0 {
		l.index++
	}

	return newToken(tokenTypeLiteral, string(l.chars[start:l.index]))
}
