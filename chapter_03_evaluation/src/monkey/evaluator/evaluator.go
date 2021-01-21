package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
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
	case *ast.PrefixExpression:
		operand := Eval(node.Operand)
		return EvalPrefixExpression(node.Operator, operand)

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

func EvalPrefixExpression(operator string, operand object.Object) object.Object {
	switch operator {
	case "!":
		return EvalBangOperatorExpression(operand)
	case "-":
		return EvalMinusOperatorExpression(operand)
	default:
		return NULL
	}
}

func EvalBangOperatorExpression(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func EvalMinusOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}