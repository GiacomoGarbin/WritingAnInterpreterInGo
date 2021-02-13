package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = EvalUnquoteCall(node, env)
	return &object.Quote{Node: node}
}

func EvalUnquoteCall(quoted ast.Node, env *object.Environment) ast.Node {
	modifier := func(node ast.Node) ast.Node {
		if !IsUnquoteCall(node) {
			return node
		}

		call, okay := node.(*ast.CallExpression)
		if !okay {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return FromObjectToASTNode(unquoted)
	}
	return ast.Modify(quoted, modifier)
}

func IsUnquoteCall(node ast.Node) bool {
	call, okay := node.(*ast.CallExpression)
	if !okay {
		return false
	}
	return call.Function.TokenLiteral() == "unquote"
}

func FromObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
        var t token.Token
        if obj.Value {
            t = token.Token{Type: token.TRUE, Literal: "true"}
        } else {
            t = token.Token{Type: token.FALSE, Literal: "false"}
        }
		return &ast.Boolean{Token: t, Value: obj.Value}
    case *object.Quote:
        return obj.Node
	default:
		return nil
	}
}