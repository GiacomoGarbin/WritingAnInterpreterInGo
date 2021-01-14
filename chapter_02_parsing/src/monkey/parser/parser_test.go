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
		value interface{}
	} {
		{"!5",     "!", 5},
		{"-15",    "-", 15},
		{"!true",  "!", true},
		{"!false", "!", false},
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
		
		if !CheckLiteralExpression(t, expression.Operand, tt.value) {
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

func CheckIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	ident, okay := expression.(*ast.Identifier)

	if !okay {
		t.Errorf("expression not *ast.Identifier, got=%T", expression)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s, got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func CheckBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, okay := expression.(*ast.Boolean)

	if !okay {
		t.Errorf("expression not *ast.Boolean, got=%T", expression)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t, got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() not %t, got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func CheckLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int:
		return CheckIntegerLiteral(t, expression, int64(value))
	case int64:
		return CheckIntegerLiteral(t, expression, value)
	case string:
		return CheckIdentifier(t, expression, value)
	case bool:
		return CheckBooleanLiteral(t, expression, value)
	}

	t.Errorf("type of expression not handled, got=%T", expression)
	return false
}

func CheckInfixExpression(
	t *testing.T,
	expression ast.Expression,
	OperandLeft interface{},
	operator string,
	OperandRight interface{},
) bool {
	operation, okay := expression.(*ast.InfixExpression)

	if !okay {
		t.Errorf("expression not *ast.InfixExpression, got=%T", expression)
		return false
	}

	if !CheckLiteralExpression(t, operation.OperandLeft, OperandLeft) {
		return false
	}

	if operation.Operator != operator {
		t.Errorf("expression.Operator not %s, got=%s", operator, operation.Operator)
		return false
	}

	if !CheckLiteralExpression(t, operation.OperandRight, OperandRight) {
		return false
	}

	return true
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input string
		OperandLeft interface{}
		operator string
		OperandRight interface{}
	} {
		{"5 + 5",  5, "+",  5},
		{"5 - 5",  5, "-",  5},
		{"5 * 5",  5, "*",  5},
		{"5 / 5",  5, "/",  5},
		{"5 < 5",  5, "<",  5},
		{"5 > 5",  5, ">",  5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !CheckInfixExpression(t, stmt.Expression, tt.OperandLeft, tt.operator, tt.OperandRight) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input string
		expected string
	} {
		{"-a + b",  "((-a) + b)"},
		{"!-a",  "(!(-a))"},
		{"a + b + c",  "((a + b) + c)"},
		{"a + b - c",  "((a + b) - c)"},
		{"a * b * c",  "((a * b) * c)"},
		{"a * b / c",  "((a * b) / c)"},
		{"a + b / c",  "(a + (b / c))"},
		{"a + b * c + d / e - f",  "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)	
		program := p.ParseProgram()
		CheckParseErrors(t, p)

		if program.String() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, program.String())
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		ExpectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
		}

		stmt, okay := program.Statements[0].(*ast.ExpressionStatement)
		if !okay {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		boolean, okay := stmt.Expression.(*ast.Boolean)
		if !okay {
			t.Fatalf("exp not *ast.Boolean, got=%T", stmt.Expression)
		}

		if boolean.Value != tt.ExpectedBoolean {
			t.Errorf("boolean.Value not %t, got=%t", tt.ExpectedBoolean, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	stmt, okay := program.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	expression, okay := stmt.Expression.(*ast.IfExpression)
	if !okay {
		t.Fatalf("expression not *ast.IfExpression, got=%T", stmt.Expression)
	}

	if !CheckInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence has more than 1 statement, got=%d", len(expression.Consequence.Statements))
	}

	consequence, okay := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("consequence is not ast.ExpressionStatement, got=%T", expression.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("alternative not nil, got=%+v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	stmt, okay := program.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	expression, okay := stmt.Expression.(*ast.IfExpression)
	if !okay {
		t.Fatalf("expression not *ast.IfExpression, got=%T", stmt.Expression)
	}

	if !CheckInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence has more than 1 statement, got=%d", len(expression.Consequence.Statements))
	}

	consequence, okay := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("consequence is not ast.ExpressionStatement, got=%T", expression.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("alternative has more than 1 statement, got=%d", len(expression.Alternative.Statements))
	}

	alternative, okay := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("alternative is not ast.ExpressionStatement, got=%T", expression.Alternative.Statements[0])
	}

	if !CheckIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	CheckParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	stmt, okay := program.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	function, okay := stmt.Expression.(*ast.FunctionLiteral)
	if !okay {
		t.Fatalf("expression not *ast.FunctionLiteral, got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.parameters has not 2 parameters, got=%d", len(function.Parameters))
	}

	CheckLiteralExpression(t, function.Parameters[0], "x")
	CheckLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.body.statements has not 1 statement, got=%d", len(function.Body.Statements))
	}

	statement, okay := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !okay {
		t.Fatalf("function.body.statements[0] is not ast.ExpressionStatement, got=%T", function.Body.Statements[0])
	}
	
	CheckInfixExpression(t, statement.Expression, "x", "+", "y")
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		ExpectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		CheckParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
		}

		stmt, okay := program.Statements[0].(*ast.ExpressionStatement)
		if !okay {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		function, okay := stmt.Expression.(*ast.FunctionLiteral)
		if !okay {
			t.Fatalf("exp not *ast.FunctionLiteral, got=%T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.ExpectedParams) {
			t.Errorf("length parameters not %d, got=%d", len(function.Parameters), len(tt.ExpectedParams))
		}

		for i, identifier := range tt.ExpectedParams {
			CheckLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}