package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func DefineMacro(program *ast.Program, env *object.Environment) {
    definitions := []int{}

    for i, statement := range program.Statements {
        if IsMacroDefinition(statement) {
            AddMacro(statement, env)
            definitions = append(definitions, i)
        }
    }

    for i := len(definitions) - 1; i >= 0; i = i - 1 {
        index := definitions[i]
        program.Statements = append(
            program.Statements[:index],
            program.Statements[index+1:]...,
        )
    }
}

func IsMacroDefinition(node ast.Statement) bool {
    stmt, okay := node.(*ast.LetStatement)
    if !okay {
        return false
    }

    _, okay = stmt.Value.(*ast.MacroLiteral)
    if !okay {
        return false
    }

    return true
}

func AddMacro(stmt ast.Statement, env *object.Environment) {
    let, _ := stmt.(*ast.LetStatement)
    literal, _ := let.Value.(*ast.MacroLiteral)

    macro := &object.Macro{
        Parameters: literal.Parameters,
        Env:        env,
        Body:       literal.Body,
    }

    env.Set(let.Name.Value, macro)
}

func ExpandMacro(program ast.Node, env *object.Environment) ast.Node {
    return ast.Modify(program, func(node ast.Node) ast.Node {
        call, okay := node.(*ast.CallExpression)
        if !okay {
            return node
        }

        macro, okay := IsMacroCall(call, env)
        if !okay {
            return node
        }

        args := QuoteArgs(call)
        env := ExtendMacroEnv(macro, args)

        evaluated := Eval(macro.Body, env)

        quote, okay := evaluated.(*object.Quote)
        if !okay {
            panic("we only support returning AST-nodes from macros")
        }

        return quote.Node
    })
}

func IsMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
    identifier, okay := exp.Function.(*ast.Identifier)
    if !okay {
        return nil, false
    }

    obj, okay := env.Get(identifier.Value)
    if !okay {
        return nil, false
    }

    macro, okay := obj.(*object.Macro)
    if !okay {
        return nil, false
    }

    return macro, true
}

func QuoteArgs(exp *ast.CallExpression) []*object.Quote {
    args := []*object.Quote{}

    for _, arg := range exp.Arguments {
        args = append(args, &object.Quote{Node: arg})
    }

    return args
}

func ExtendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
    extended := object.NewEnclosedEnvironment(macro.Env)

    for index, param := range macro.Parameters {
        extended.Set(param.Value, args[index])
    }

    return extended
}