package lexer

import "monkey/token"

type Lexer struct {
	input string
	position int
	ReadPosition int
	char byte
}

func (lexer *Lexer) ReadChar() {
	if lexer.ReadPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.ReadPosition]
	}
	lexer.position = lexer.ReadPosition
	lexer.ReadPosition += 1
}

func (lexer *Lexer) PeekChar() byte {
	if lexer.ReadPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.ReadPosition]
	}
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.ReadChar()
	return lexer
}

func NewToken(TokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: TokenType, Literal: string(char)}
}

func IsLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func (lexer *Lexer) ReadIdentifier() string {
	position := lexer.position
	for IsLetter(lexer.char) {
		lexer.ReadChar()
	}
	return lexer.input[position:lexer.position]
}

func IsDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lexer *Lexer) ReadNumber() string {
	position := lexer.position
	for IsDigit(lexer.char) {
		lexer.ReadChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) SkipWhiteSpace() {
	for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\n' || lexer.char == '\r' {
		lexer.ReadChar()
	}
}

func (lexer *Lexer) NextToken() token.Token {
	var t token.Token

	lexer.SkipWhiteSpace()

	switch lexer.char {
	case '=':
		if lexer.PeekChar() == '=' {
			t.Type = token.EQ
			t.Literal = "=="
			lexer.ReadChar()
		} else {
			t = NewToken(token.ASSIGN, lexer.char)
		}
	case ';':
		t = NewToken(token.SEMICOLON, lexer.char)
	case '(':
		t = NewToken(token.LPAREN, lexer.char)
	case ')':
		t = NewToken(token.RPAREN, lexer.char)
	case ',':
		t = NewToken(token.COMMA, lexer.char)
	case '+':
		t = NewToken(token.PLUS, lexer.char)
	case '-':
		t = NewToken(token.MINUS, lexer.char)
	case '!':
		if lexer.PeekChar() == '=' {
			t.Type = token.NOT_EQ
			t.Literal = "!="
			lexer.ReadChar()
		} else {
			t = NewToken(token.BANG, lexer.char)
		}
	case '*':
		t = NewToken(token.ASTERISK, lexer.char)
	case '/':
		t = NewToken(token.SLASH, lexer.char)
	case '<':
		t = NewToken(token.LT, lexer.char)
	case '>':
		t = NewToken(token.GT, lexer.char)
	case '{':
		t = NewToken(token.LBRACE, lexer.char)
	case '}':
		t = NewToken(token.RBRACE, lexer.char)
	case 0:
		t.Type = token.EOF
		t.Literal = ""
	default:
		if IsLetter(lexer.char) {
			t.Literal = lexer.ReadIdentifier()
			t.Type = token.LookUpIdent(t.Literal)
			return t
		} else if IsDigit(lexer.char) {
			t.Literal = lexer.ReadNumber()
			t.Type = token.INT
			return t
		} else {
			t = NewToken(token.ILLEGAL, lexer.char)
		}
	}

	lexer.ReadChar()
	return t
}