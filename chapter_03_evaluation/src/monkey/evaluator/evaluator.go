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
		return EvalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		operand := Eval(node.Operand)
		return EvalPrefixExpression(node.Operator, operand)
	case *ast.InfixExpression:
		OperandLeft  := Eval(node.OperandLeft)
		OperandRight := Eval(node.OperandRight)
		return EvalInfixExpression(node.Operator, OperandLeft, OperandRight)
	case *ast.BlockStatement:
		return EvalBlockStatement(node)
	case *ast.ReturnStatement:
		value := Eval(node.Value)
		return &object.ReturnValue{Value: value}

	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return BoolToBoolean(node.Value)
	case *ast.IfExpression:
		return EvalIfElseExpression(node)
	}

	return nil
}

func EvalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		if ReturnValue, okay := result.(*object.ReturnValue); okay {
			return ReturnValue.Value
		}
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

func EvalInfixExpression(operator string, OperandLeft object.Object, OperandRight object.Object) object.Object {
	switch {
	case OperandLeft.Type() == object.INTEGER_OBJ && OperandRight.Type() == object.INTEGER_OBJ:
		return EvalIntegerInfixExpression(operator, OperandLeft, OperandRight)
	case operator == "==":
		return BoolToBoolean(OperandLeft == OperandRight)
	case operator == "!=":
		return BoolToBoolean(OperandLeft != OperandRight)
	default:
		return NULL
	}
}

func EvalIntegerInfixExpression(operator string, OperandLeft object.Object, OperandRight object.Object) object.Object {
	ValueLeft  := OperandLeft.(*object.Integer).Value
	ValueRight := OperandRight.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: ValueLeft + ValueRight}
	case "-":
		return &object.Integer{Value: ValueLeft - ValueRight}
	case "*":
		return &object.Integer{Value: ValueLeft * ValueRight}
	case "/":
		return &object.Integer{Value: ValueLeft / ValueRight}
	case "<":
		return BoolToBoolean(ValueLeft < ValueRight)
	case ">":
		return BoolToBoolean(ValueLeft > ValueRight)
	case "==":
		return BoolToBoolean(ValueLeft == ValueRight)
	case "!=":
		return BoolToBoolean(ValueLeft != ValueRight)
	default:
		return NULL
	}
}

func EvalIfElseExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if IsTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func IsTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func EvalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}