package mirrormeta

import "testing"

func Test_Lex(t *testing.T) {
	type test struct {
		input    string
		expected []token
	}

	tests := []test{
		{
			"type:{ name?: string, email?: string }",
			[]token{
				newToken(tokenTypeLiteral, "type"),
				newToken(tokenTypeColon, ":"),
				newToken(tokenTypeLiteral, "{ name?: string, email?: string }"),
			},
		},
		{
			"type:string, name:email, optional:true",
			[]token{
				newToken(tokenTypeLiteral, "type"),
				newToken(tokenTypeColon, ":"),
				newToken(tokenTypeLiteral, "string"),
				newToken(tokenTypeComma, ","),
				newToken(tokenTypeLiteral, "name"),
				newToken(tokenTypeColon, ":"),
				newToken(tokenTypeLiteral, "email"),
				newToken(tokenTypeComma, ","),
				newToken(tokenTypeLiteral, "optional"),
				newToken(tokenTypeColon, ":"),
				newToken(tokenTypeLiteral, "true"),
			},
		},
	}

	for _, test := range tests {
		lexer := newLexer(test.input)
		for _, expected := range test.expected {
			actual := lexer.lex()
			if actual != expected {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	}
}
