package parser

import (
	"fmt"
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

	// for i, stmt := range program.Statements {
	// 	t.Errorf("program.Statements[%d] = %T (%s, %s)", i, stmt, stmt.TokenLiteral(), stmt.String())
	// }

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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

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

	literal, okay := stmt.Expression.(*ast.IntegerLiteral)

	if !okay {
		t.Fatalf("stmt.Expression is not ast.Identifier, got=%T",  stmt.Expression)
	}
	
	if literal.Value != 5 {
		t.Errorf("literal.Value not '%d', got=%d", 5, literal.Value)
	}
	
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not '%s', got=%s", "5", literal.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input string
		operator string
		value int64
	} {
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
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

		expression, okay := stmt.Expression.(*ast.PrefixExpression)
	
		if !okay {
			t.Fatalf("stmt is not ast.PrefixExpression, got=%T",  stmt.Expression)
		}
		
		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not '%s', got=%s", tt.operator, expression.Operator)
		}
		
		if !CheckIntegerLiteral(t, expression.Operand, tt.value) {
			return
		}
	}
}

func CheckIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	literal, okay := il.(*ast.IntegerLiteral)

	if !okay {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return false
	}

	if literal.Value != value {
		t.Errorf("literal.Value not %d, got=%d", value, literal.Value)
		return false
	}

	if literal.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("literal.TokenLiteral() not %d, got=%s", value, literal.TokenLiteral())
		return false
	}

	return true
}