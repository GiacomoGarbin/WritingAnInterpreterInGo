package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {

	switch node := node.(type) {

	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)

	case *InfixExpression:
		node.OperandLeft, _ = Modify(node.OperandLeft, modifier).(Expression)
		node.OperandRight, _ = Modify(node.OperandRight, modifier).(Expression)

	case *PrefixExpression:
		node.Operand, _ = Modify(node.Operand, modifier).(Expression)

	case *IndexExpression:
		node.Array, _ = Modify(node.Array, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)

	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}

	case *BlockStatement:
		for i, _ := range node.Statements {
			node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
		}

	case *ReturnStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *FunctionLiteral:
		for i, _ := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	case *ArrayLiteral:
		for i, _ := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
		}

    case *HashLiteral:
        pairs := make(map[Expression]Expression)
        for OldKey, OldVal := range node.Pairs {
            NewKey, _ := Modify(OldKey, modifier).(Expression)
            NewVal, _ := Modify(OldVal, modifier).(Expression)
            pairs[NewKey] = NewVal
        }
        node.Pairs = pairs

	}

	return modifier(node)
}