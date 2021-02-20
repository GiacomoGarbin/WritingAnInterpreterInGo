package evaluator

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func CheckParseProgram(input string) *ast.Program {
    l := lexer.NewLexer(input)
    p := parser.NewParser(l)
    return p.ParseProgram()
}

func TestDefineMacro(t *testing.T) {
    input := `let number = 1;
    let function = fn(x, y) { x + y; };
    let mymacro = macro(x, y) { x + y; };`

    env := object.NewEnvironment()
    program := CheckParseProgram(input)

    DefineMacro(program, env)

    if len(program.Statements) != 2 {
        t.Fatalf("wrong number of statements, got=%d", len(program.Statements))
    }

    _, okay := env.Get("number")
    if okay {
        t.Fatalf("number should not be defined")
	}
	
    _, okay = env.Get("function")
    if okay {
        t.Fatalf("function should not be defined")
    }

    obj, okay := env.Get("mymacro")
    if !okay {
        t.Fatalf("macro not in environment")
    }

    macro, okay := obj.(*object.Macro)
    if !okay {
        t.Fatalf("object is not macro, got=%T (%+v)", obj, obj)
    }

    if len(macro.Parameters) != 2 {
        t.Fatalf("wrong number of macro parameters, got=%d", len(macro.Parameters))
    }

    if macro.Parameters[0].String() != "x" {
        t.Fatalf("parameter is not 'x', got=%q", macro.Parameters[0])
    }
    if macro.Parameters[1].String() != "y" {
        t.Fatalf("parameter is not 'y', got=%q", macro.Parameters[1])
    }

    body := "(x + y)"
    if macro.Body.String() != body {
        t.Fatalf("body is not %q, got=%q", body, macro.Body.String())
    }
}

func TestExpandMacro(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {
            `let infixExpression = macro() { quote(1 + 2); };
            infixExpression();`,
            `(1 + 2)`,
        },
        {
            `let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
            reverse(2 + 2, 10 - 5);`,
            `(10 - 5) - (2 + 2)`,
        },
        {
            `let unless = macro(condition, consequence, alternative) {
                quote(if(!(unquote(condition))) {
                    unquote(consequence);
                } else {
                    unquote(alternative);
                });
            };
            unless(10 > 5, puts("not greater"), puts("greater"));`,
            `if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
        },
    }

    for _, tt := range tests {
        expected := CheckParseProgram(tt.expected)
        program := CheckParseProgram(tt.input)

        env := object.NewEnvironment()
        DefineMacro(program, env)
        expanded := ExpandMacro(program, env)

        if expanded.String() != expected.String() {
            t.Errorf("not equal, want=%q, got=%q", expected.String(), expanded.String())
        }
    }
}