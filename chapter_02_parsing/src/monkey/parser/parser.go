package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer *lexer.Lexer

	CurrToken token.Token
	PeekToken token.Token

	errors []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	// set both CurrToken and PeekToken
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
	p.CurrToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

func (p *Parser) CurrTokenIs(t token.TokenType) bool {
	return p.CurrToken.Type == t
}

func (p *Parser) PeekTokenIs(t token.TokenType) bool {
	return p.PeekToken.Type == t
}

func (p *Parser) ExpectedPeek(t token.TokenType) bool {
	if p.PeekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.AddError(t)
		return false
	}
}

func (p *Parser) GetErrors() []string {
	return p.errors
}

func (p *Parser) AddError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s insted", t, p.PeekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.CurrToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.CurrToken}


	fmt.Printf("stuck!")

	if !p.ExpectedPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.CurrToken, Value: p.CurrToken.Literal};
	
	if !p.ExpectedPeek(token.ASSIGN) {
		return nil
	}

	if !p.CurrTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	
	for p.CurrToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	
	return program
}