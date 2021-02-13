package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	TurnOneIntoTwo := func(node Node) Node {
		integer, okay := node.(*IntegerLiteral)

		if !okay {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: one()},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: two()},
				},
			},
		},
        {
            &InfixExpression{OperandLeft: one(), Operator: "+", OperandRight: two()},
            &InfixExpression{OperandLeft: two(), Operator: "+", OperandRight: two()},
        },
        {
            &InfixExpression{OperandLeft: two(), Operator: "+", OperandRight: one()},
            &InfixExpression{OperandLeft: two(), Operator: "+", OperandRight: two()},
        },
        {
            &PrefixExpression{Operator: "-", Operand: one()},
            &PrefixExpression{Operator: "-", Operand: two()},
        },
        {
            &IndexExpression{Array: one(), Index: one()},
            &IndexExpression{Array: two(), Index: two()},
		},
		{
			&IfExpression{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&IfExpression{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
        {
            &ReturnStatement{Value: one()},
            &ReturnStatement{Value: two()},
        },
        {
            &LetStatement{Value: one()},
            &LetStatement{Value: two()},
        },
        {
            &FunctionLiteral{
                Parameters: []*Identifier{},
                Body: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: one()},
                    },
                },
            },
            &FunctionLiteral{
                Parameters: []*Identifier{},
                Body: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: two()},
                    },
                },
            },
        },
        {
            &ArrayLiteral{Elements: []Expression{one(), one()}},
            &ArrayLiteral{Elements: []Expression{two(), two()}},
        },
	}

	for _, tt := range tests {
		modified := Modify(tt.input, TurnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal, got=%#v, want=%#v", modified, tt.expected)
		}
	}

    hash := &HashLiteral{
        Pairs: map[Expression]Expression{
            one(): one(),
            one(): one(),
        },
    }

    Modify(hash, TurnOneIntoTwo)

    for key, val := range hash.Pairs {
        key, _ := key.(*IntegerLiteral)
        if key.Value != 2 {
            t.Errorf("value is not %d, got=%d", 2, key.Value)
        }
        val, _ := val.(*IntegerLiteral)
        if val.Value != 2 {
            t.Errorf("value is not %d, got=%d", 2, val.Value)
        }
    }
}