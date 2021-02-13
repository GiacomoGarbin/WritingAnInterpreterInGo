package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return EvalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		operand := Eval(node.Operand, env)
		if IsError(operand) {
			return operand
		}
		return EvalPrefixExpression(node.Operator, operand)
	case *ast.InfixExpression:
		OperandLeft  := Eval(node.OperandLeft, env)
		if IsError(OperandLeft) {
			return OperandLeft
		}
		OperandRight := Eval(node.OperandRight, env)
		if IsError(OperandRight) {
			return OperandRight
		}
		return EvalInfixExpression(node.Operator, OperandLeft, OperandRight)
	case *ast.BlockStatement:
		return EvalBlockStatement(node, env)
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if IsError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
	case *ast.ReturnStatement:
		value := Eval(node.Value, env)
		if IsError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}

	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return BoolToBoolean(node.Value)
	case *ast.IfExpression:
		return EvalIfElseExpression(node, env)

	case *ast.Identifier:
		return EvalIdentifier(node, env)
		
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Env: env,
			Body: node.Body}

	case *ast.CallExpression:
		if node.Function.TokenLiteral() == "quote" {
			return quote(node.Arguments[0], env)
		}

		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}

		args := EvalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}

		return CallFunction(function, args)
	
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := EvalExpressions(node.Elements, env)

		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		array := Eval(node.Array, env)
		if IsError(array) {
			return array
		}
		
		index := Eval(node.Index, env)
		if IsError(index) {
			return index
		}

		return EvalIndexExpression(array, index)

	case *ast.HashLiteral:
		return EvalHashLiteral(node, env)
	}

	return nil
}

func EvalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
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
		return NewError("unknown operator: %s%s", operator, operand.Type())
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
		return NewError("unknown operator: -%s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func EvalInfixExpression(operator string, OperandLeft object.Object, OperandRight object.Object) object.Object {
	switch {
	case OperandLeft.Type() == object.INTEGER_OBJ && OperandRight.Type() == object.INTEGER_OBJ:
		return EvalIntegerInfixExpression(operator, OperandLeft, OperandRight)
	case OperandLeft.Type() == object.STRING_OBJ && OperandRight.Type() == object.STRING_OBJ:
		return EvalStringInfixExpression(operator, OperandLeft, OperandRight)
	case operator == "==":
		return BoolToBoolean(OperandLeft == OperandRight)
	case operator == "!=":
		return BoolToBoolean(OperandLeft != OperandRight)
	case OperandLeft.Type() != OperandRight.Type():
		return NewError("type mismatch: %s %s %s", OperandLeft.Type(), operator, OperandRight.Type())
	default:
		return NewError("unknown operator: %s %s %s", OperandLeft.Type(), operator, OperandRight.Type())
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
		return NewError("unknown operator: %s %s %s", OperandLeft.Type(), operator, OperandRight.Type())
	}
}

func EvalStringInfixExpression(operator string, OperandLeft object.Object, OperandRight object.Object) object.Object {
	ValueLeft  := OperandLeft.(*object.String).Value
	ValueRight := OperandRight.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: ValueLeft + ValueRight}

	default:
		return NewError("unknown operator: %s %s %s", OperandLeft.Type(), operator, OperandRight.Type())
	}
}

func EvalIfElseExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if IsError(condition) {
		return condition
	}

	if IsTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
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

func EvalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			ResultType := result.Type()

			if ResultType == object.RETURN_VALUE_OBJ || ResultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func NewError(format string, a ... interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func EvalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	
	if value, okay := env.Get(i.Value); okay {
		return value
	}

	if builtin, okay := builtins[i.Value]; okay {
		return builtin
	}

	return NewError("identifier not found: " + i.Value)
}

func EvalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Eval(expression, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func CallFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		ExtendedEnv := ExtendFunctionEnv(function, args)
		evaluated := Eval(function.Body, ExtendedEnv)
		return UnwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Func(args...)
	default:
		return NewError("not a function: %s", fn.Type())
	}
}

func ExtendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func UnwrapReturnValue(obj object.Object) object.Object {
	if ReturnValue, okay := obj.(*object.ReturnValue); okay {
		return ReturnValue.Value
	}
	return obj
}

func EvalIndexExpression(container, index object.Object) object.Object {
	switch {
	case container.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return EvalArrayIndexExpression(container, index)
	case container.Type() == object.HASH_OBJ:
		return EvalHashIndexExpression(container, index)
	default:
		return NewError("index operator not supported: %s", container.Type())
	}
}

func EvalArrayIndexExpression(container, index object.Object) object.Object {
	array := container.(*object.Array)
	i := index.(*object.Integer).Value
	size := int64(len(array.Elements))

	if i < 0 || i >= size {
		return NULL
	}

	return array.Elements[i]
}

func EvalHashIndexExpression(container, index object.Object) object.Object {
	hash := container.(*object.Hash)

	key, okay := index.(object.Hashable)
	if !okay {
		return NewError("unusable as hash key: %s", index.Type())
	}

	pair, okay := hash.Pairs[key.HashKey()]
	if !okay {
		return NULL
	}

	return pair.Value
}

func EvalHashLiteral(hash *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for k, v := range hash.Pairs {
		key := Eval(k, env)
		if IsError(key) {
			return key
		}

		hashable, okay := key.(object.Hashable)
		if !okay {
			return NewError("unusable as hash key: %s", key.Type())
		}

		val := Eval(v, env)
		if IsError(val) {
			return val
		}

		pairs[hashable.HashKey()] = object.HashPair{Key: key, Value: val}
	}

	return &object.Hash{Pairs: pairs}
}