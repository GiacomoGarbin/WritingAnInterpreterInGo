package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return EvalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return BoolToBoolean(node.Value)
	}

	return nil
}

func EvalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func BoolToBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}