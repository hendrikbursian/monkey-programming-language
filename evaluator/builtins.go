package evaluator

import (
	"fmt"

	"github.com/hendrikbursian/monkey-programming-language/object"
)

var builtins map[string]*object.Builtin = map[string]*object.Builtin{
	"l": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}, nil
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}, nil
			default:
				return nil, fmt.Errorf("argument to `l` not supported. got=%s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("wrong number of arguments to push. got=%d, want=%d", len(args), 2)
			}

			arrObj, ok := args[0].(*object.Array)
			if !ok {
				return nil, fmt.Errorf("first argument to push has to be an array, got %s instead", args[0].Type())
			}

			arrObj.Elements = append(arrObj.Elements, args[1])

			return arrObj, nil
		},
	},
	"first": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments to first, got=%d, want=%d", len(args), 1)
			}

			arrObj, ok := args[0].(*object.Array)
			if !ok {
				return nil, fmt.Errorf("first argument to first has to be an array, got %s instead", args[0].Type())

			}

			// TODO: implement maybes
			return arrObj.Elements[0], nil
		},
	},
	"last": {
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments to last, got=%d, want=%d", len(args), 1)
			}

			arrObj, ok := args[0].(*object.Array)
			if !ok {
				return nil, fmt.Errorf("first argument to last has to be an array, got %s instead", args[0].Type())

			}

			// TODO implement maybes
			return arrObj.Elements[len(arrObj.Elements)-1], nil
		},
	},
}
