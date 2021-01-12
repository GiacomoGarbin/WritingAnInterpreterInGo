package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS		// ==
	LESSGREATER	// < or >
	SUM			// +
	PRODUCT		// *
	PREFIX		// - or !
	CALL		// func(args)
)

var precedences = map[token.TokenType] int {
	token.EQ:		EQUALS,
	token.NOT_EQ:	EQUALS,
	token.LT:		LESSGREATER,
	token.GT:		LESSGREATER,
	token.PLUS:		SUM,
	token.MINUS:	SUM,
	token.SLASH:	PRODUCT,
	token.ASTERISK:	PRODUCT,
}

type (
	PrefixParseFn func() ast.Expression
	InfixParseFn func(ast.Expression) ast.Expression
)

func (p *Parser) RegisterPrefixParseFn(TokenType token.TokenType, fn PrefixParseFn) {
	p.PrefixParseFns[TokenType] = fn
}

func (p *Parser) RegisterInfixParseFn(TokenType token.TokenType, fn InfixParseFn) {
	p.InfixParseFns[TokenType] = fn
}

type Parser struct {
	lexer *lexer.Lexer

	CurrToken token.Token
	PeekToken token.Token

	errors []string

	PrefixParseFns map[token.TokenType] PrefixParseFn
	InfixParseFns map[token.TokenType] InfixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	// init prefix parse functions map
	p.PrefixParseFns = make(map[token.TokenType] PrefixParseFn)
	// register prefix parsing functions
	p.RegisterPrefixParseFn(token.IDENT, p.ParseIdentifier)
	p.RegisterPrefixParseFn(token.INT, p.ParseIntegerLiteral)
	p.RegisterPrefixParseFn(token.BANG, p.ParsePrefixExpression)
	p.RegisterPrefixParseFn(token.MINUS, p.ParsePrefixExpression)
	p.RegisterPrefixParseFn(token.TRUE, p.ParseBoolean)
	p.RegisterPrefixParseFn(token.FALSE, p.ParseBoolean)
	
	// init infix parse functions map
	p.InfixParseFns = make(map[token.TokenType] InfixParseFn)
	// register infix parsing functions
	p.RegisterInfixParseFn(token.PLUS, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.MINUS, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.SLASH, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.ASTERISK, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.EQ, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.NOT_EQ, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.LT, p.ParseInfixExpression)
	p.RegisterInfixParseFn(token.GT, p.ParseInfixExpression)

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

func (p *Parser) CurrPrecedence() int {
	if precedence, okay := precedences[p.CurrToken.Type]; okay {
		return precedence
	}
	return LOWEST
}

func (p *Parser) PeekPrecedence() int {
	if precedence, okay := precedences[p.PeekToken.Type]; okay {
		return precedence
	}
	return LOWEST
}

func (p *Parser) ExpectedPeek(t token.TokenType) bool {
	if p.PeekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.ExpectedPeekError(t)
		return false
	}
}

func (p *Parser) GetErrors() []string {
	return p.errors
}

func (p *Parser) ExpectedPeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s insted", t, p.PeekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NoPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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

func (p *Parser) ParseStatement() ast.Statement {
	switch p.CurrToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.CurrToken}

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

	p.NextToken() // we shouldn't need this

	return stmt
}

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.CurrToken}

	p.NextToken()

	if !p.CurrTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.CurrToken}

	stmt.Expression = p.ParseExpression(LOWEST)

	if p.PeekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := p.PrefixParseFns[p.CurrToken.Type]

	if prefix == nil {
		p.NoPrefixParseFnError(p.CurrToken.Type)
		return nil
	}

	LeftExpression := prefix()

	for !p.PeekTokenIs(token.SEMICOLON) && precedence < p.PeekPrecedence() {
		infix := p.InfixParseFns[p.PeekToken.Type]

		if infix == nil {
			return LeftExpression
		}

		p.NextToken()

		LeftExpression = infix(LeftExpression)
	}

	return LeftExpression
}

func (p *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.CurrToken, Value: p.CurrToken.Literal}
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.CurrToken}

	value, err := strconv.ParseInt(p.CurrToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.CurrToken.Literal)
		p.errors = append(p.errors, msg)
	}

	literal.Value = value

	return literal
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.CurrToken, Operator: p.CurrToken.Literal}

	p.NextToken()

	expression.Operand = p.ParseExpression(PREFIX)

	return expression
}

func (p *Parser) ParseInfixExpression(operand ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.CurrToken,
		OperandLeft: operand,
		Operator: p.CurrToken.Literal,
	}

	precedence := p.CurrPrecedence()
	p.NextToken()
	expression.OperandRight = p.ParseExpression(precedence)

	return expression
}

func (p *Parser) ParseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.CurrToken, Value: p.CurrTokenIs(token.TRUE)}
}