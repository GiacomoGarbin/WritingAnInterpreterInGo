package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
		x + y;
	};

	let result = add(five, ten);
	`

	tests := [] struct {
		ExpectedType	token.TokenType
		ExpectedLiteral string
	} {
		{ token.LET, "let" },
		{ token.IDENT, "five" },
		{ token.ASSIGN, "=" },
		{ token.INT, "5" },
		{ token.SEMICOLON, ";" },
		{ token.LET, "let" },
		{ token.IDENT, "ten" },
		{ token.ASSIGN, "=" },
		{ token.INT, "10" },
		{ token.SEMICOLON, ";" },
		{ token.LET, "let" },
		{ token.IDENT, "add" },
		{ token.ASSIGN, "=" },
		{ token.FUNCTION, "fn" },
		{ token.LPAREN, "(" },
		{ token.IDENT, "x" },
		{ token.COMMA, "," },
		{ token.IDENT, "y" },
		{ token.RPAREN, ")" },
		{ token.LBRACE, "{" },
		{ token.IDENT, "x" },
		{ token.PLUS, "+" },
		{ token.IDENT, "y" },
		{ token.SEMICOLON, ";" },
		{ token.RBRACE, "}" },
		{ token.SEMICOLON, ";" },
		{ token.LET, "let" },
		{ token.IDENT, "result" },
		{ token.ASSIGN, "=" },
		{ token.IDENT, "add" },
		{ token.LPAREN, "(" },
		{ token.IDENT, "five" },
		{ token.COMMA, "," },
		{ token.IDENT, "ten" },
		{ token.RPAREN, ")" },
		{ token.SEMICOLON, ";" },
		{ token.EOF, "" },
	}

	lexer := NewLexer(input)

	for i, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.ExpectedType {
			t.Fatalf("tests[%d] token type wrong, expected=%q, got=%q", i, test.ExpectedType, token.Type)
		}

		if token.Literal != test.ExpectedLiteral {
			t.Fatalf("tests[%d] token literal wrong, expected=%q, got=%q", i, test.ExpectedLiteral, token.Literal)
		}
	}
}