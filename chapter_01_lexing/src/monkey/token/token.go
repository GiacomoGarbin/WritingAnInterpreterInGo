package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	INT = "INT"

	ASSIGN = "="
	PLUS = "+"

	COMMA = ","
	SEMICOLON = ";"
	
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET = "LET"
)

var keywords = map[string] TokenType {
	"fn": FUNCTION,
	"let": LET,
}

func LookUpIdent(ident string) TokenType {
	if token, okay := keywords[ident]; okay {
		return token
	}
	return IDENT
}