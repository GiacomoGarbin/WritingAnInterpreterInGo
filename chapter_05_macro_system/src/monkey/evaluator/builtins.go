package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin {
	"len": &object.Builtin{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments, got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return NewError("argument to len not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments, got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("builtin first argument must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments, got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("builtin last argument must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			size := len(array.Elements)
			if size > 0 {
				return array.Elements[size - 1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments, got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("builtin rest argument must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			size := len(array.Elements)
			if size > 0 {
				elements := make([]object.Object, size - 1, size - 1)
				copy(elements, array.Elements[1:size])
				return &object.Array{Elements: elements}
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("wrong number of arguments, got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("builtin push first argument must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			size := len(array.Elements)

			elements := make([]object.Object, size + 1, size + 1)
			copy(elements, array.Elements[0:size])
			elements[size] = args[1]

			return &object.Array{Elements: elements}
		},
	},
	"puts": &object.Builtin{
		Func: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}