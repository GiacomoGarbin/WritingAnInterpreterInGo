package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func LetStatementTest(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let', got=%q", s.TokenLiteral())
		return false
	}

	LetStmt, okay := s.(*ast.LetStatement)
	if !okay {
		t.Errorf("s not *ast.LetStatement, got=%T", s)
		return false
	}

	if LetStmt.Name.Value != name {
		t.Errorf("LetStmt.name.value not '%s', got=%s", name, LetStmt.Name.Value)
		return false
	}

	if LetStmt.Name.TokenLiteral() != name {
		t.Errorf("LetStmt.name.TokenLiteral() not '%s', got=%s", name, LetStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func CheckParseErrors(t *testing.T, p *Parser) {
	errors := p.GetErrors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}

	t.FailNow()
}

func TestLetStatement(t *testing.T) {
	input := `let x = 5;
	let y = 10;
	let foobar = 838383;`

	// failing test input
	// input := `let x 5;
	// let = 10;
	// let 838383;`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	CheckParseErrors(t, p)

	for i, stmt := range program.Statements {
		t.Errorf("program.Statements[%d] = %T (%s, %s)", i, stmt, stmt.TokenLiteral(), stmt.String())
	}

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	tests := []struct {
		ExpectedIdentifier string
	} {
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !LetStatementTest(t, stmt, tt.ExpectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `return 5;
	return 10;
	return 993322;`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	CheckParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		rs, okay := stmt.(*ast.ReturnStatement)
		
		if !okay {
			t.Errorf("stmt not *ast.ReturnStatement, got=%T", stmt)
			continue
		}
		
		if rs.TokenLiteral() != "return" {
			t.Errorf("rs.TokenLiteral() not 'return', got=%q", rs.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	CheckParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
	}

	stmt, okay := program.Statements[0].(*ast.ExpressionStatement)

	if !okay {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, okay := stmt.Expression.(*ast.Identifier)

	if !okay {
		t.Fatalf("stmt.Expression is not ast.Identifier, got=%T",  stmt.Expression)
	}
	
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not '%s', got=%s", "foobar", ident.Value)
	}
	
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not '%s', got=%s", "foobar", ident.TokenLiteral())
	}
}