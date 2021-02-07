package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckBooleanObject(t, evaluated, tt.expected)
	}
}

func CheckEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func CheckIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, okay := obj.(*object.Integer)

	if !okay {
		t.Errorf("object is not integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func CheckBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, okay := obj.(*object.Boolean)

	if !okay {
		t.Errorf("object is not boolean, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		integer, okay := tt.expected.(int)

		if okay {
			CheckIntegerObject(t, evaluated, int64(integer))
		} else {
			CheckNullObject(t, evaluated)
		}
	}
}

func CheckNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckIntegerObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`, 10},
		{`let f = fn(x) {
			return x;
			x + 10;
		};
		f(10);`,
		10},
		{`let f = fn(x) {
			let result = x + 10;
			return result;
			return 10;
		};
		f(10);`,
		20},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		ExpectedMessage string
	}{
		{ "5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{ "5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{ "-true;", "unknown operator: -BOOLEAN"},
		{ "true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{ "5; true + false; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{ "if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if (10 > 1) {
			if (10 > 1) {
				return true + false;
			}
			return 1;
		}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{ "foobar", "identifier not found: foobar" },
		{ "\"Hello\" - \"World!\";", "unknown operator: STRING - STRING" },
		{ `{"name": "monkey"}[fn(x) { x }];`, "unusable as hash key: FUNCTION" },
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		
		ErrorObject, okay := evaluated.(*object.Error)

		if !okay {
			t.Errorf("no error object returned, got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if ErrorObject.Message != tt.ExpectedMessage {
			t.Errorf("wrong error message, got=%q, want=%q", ErrorObject.Message, tt.ExpectedMessage)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := CheckEval(input)

	function, okay := evaluated.(*object.Function)

	if !okay {
		t.Fatalf("object is not a function, got=%T(%+v)", evaluated, evaluated)
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("function has wrong parameters, got=%+v", function.Parameters)
	}

	if function.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", function.Parameters[0].String())
	}

	if function.Body.String() != "(x + 2)" {
		t.Fatalf("body is not '(x + 2)', got=%q", function.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{ "let identity = fn(x) { x; }; identity(5);", 5 },
		{ "let identity = fn(x) { return x; }; identity(5);", 5 },
		{ "let double = fn(x) { x * 2; }; double(5);", 10 },
		{ "let add = fn(x, y) { x + y; }; add(5, 5);", 10 },
		{ "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20 },
		{ "fn(x) { x; }(5);", 5 },
	}

	for _, tt := range tests {
		CheckIntegerObject(t, CheckEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{ `let NewAdder = fn(x) { fn(y) { x + y; }; };
		let AddTwo = NewAdder(2);
		AddTwo(3);`, 5 },
		{ `let NewAdder = fn(x) { fn(y) { x + y; }; };
		let AddThree= NewAdder(3);
		AddThree(7);`, 10 },
	}

	for _, tt := range tests {
		CheckIntegerObject(t, CheckEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := "\"Hello World!\";"
	evaluated := CheckEval(input)

	str, okay := evaluated.(*object.String)

	if !okay {
		t.Fatalf("object is not a string, got=%T(%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("string has wrong value, got=%q, want=\"Hello World!\"", str.Value)
	}
}

func TestStrinConcatenation(t *testing.T) {
	input := "\"Hello\" + \" \" + \"World!\";"
	evaluated := CheckEval(input)

	str, okay := evaluated.(*object.String)

	if !okay {
		t.Fatalf("object is not a string, got=%T(%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("string has wrong value, got=%q, want=\"Hello World!\"", str.Value)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{ `len("")`, 0 },
		{ `len("four")`, 4 },
		{ `len("hello world")`, 11 },
		{ `len(1)`, "argument to len not supported, got INTEGER" },
		{ `len("one", "two")`, "wrong number of arguments, got=2, want=1" },
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			CheckIntegerObject(t, evaluated, int64(expected))
		case string:
			err, okay := evaluated.(*object.Error)

			if !okay {
				t.Errorf("object is not Error, got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if err.Message != expected {
				t.Errorf("wrong error message, expected=%q, got=%q", expected, err.Message)
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := CheckEval(input)

	array, okay := evaluated.(*object.Array)

	if !okay {
		t.Fatalf("object is not an array, got=%T(%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array has wrong number of elements, got=%d", len(array.Elements))
	}
	
	CheckIntegerObject(t, array.Elements[0], 1)
	CheckIntegerObject(t, array.Elements[1], 4)
	CheckIntegerObject(t, array.Elements[2], 6)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{ `[1, 2, 3][0];`, 1 },
		{ `[1, 2, 3][1];`, 2 },
		{ `[1, 2, 3][2];`, 3 },
		{ `let i = 0; [1][i];`, 1 },
		{ `[1, 2, 3][1 + 1];`, 3 },
		{ `let array = [1, 2, 3]; array[2]`, 3 },
		{ `let array = [1, 2, 3]; array[0] + array[1] + array[2];`, 6 },
		{ `let array = [1, 2, 3]; let i = array[0]; array[i];`, 2 },
		{ `[1, 2, 3][3];`, nil },
		{ `[1, 2, 3][-1];`, nil },
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)

		integer, okay := tt.expected.(int)

		if okay {
			CheckIntegerObject(t, evaluated, int64(integer))
		} else {
			CheckNullObject(t, evaluated)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := CheckEval(input)

	hash, okay := evaluated.(*object.Hash)
	if !okay {
		t.Fatalf("eval didn't return hash, got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey] int64 {
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(hash.Pairs) != len(expected) {
		t.Fatalf("hash has wrong number of pairs, got=%d", len(hash.Pairs))
	}

	for key, val := range expected {
		pair, okay := hash.Pairs[key]
		if !okay {
			t.Errorf("no pair for key in pairs")
		}
		CheckIntegerObject(t, pair.Value, val)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{ `{"foo": 5}["foo"];`, 5 },
		{ `{"foo": 5}["bar"];`, nil },
		{ `let key = "foo"; {"foo": 5}[key];`, 5 },
		{ `{}["foo"];`, nil },
		{ `{5: 5}[5];`, 5 },
		{ `{true: 5}[true];`, 5 },
		{ `{false: 5}[false];`, 5 },
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		integer, okay := tt.expected.(int)
		if okay {
			CheckIntegerObject(t, evaluated, int64(integer))
		} else {
			CheckNullObject(t, evaluated)
		}
	}
}