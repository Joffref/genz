package parser

import (
	"go/ast"
	"go/types"
)

func parseMethods(signatures map[string]*types.Signature) ([]Method, error) {
	methods := make([]Method, len(signatures))

	index := 0
	for name, signature := range signatures {

		params := []Type{}
		if signature.Params() != nil {
			params = make([]Type, signature.Params().Len())
			for j := 0; j < signature.Params().Len(); j++ {
				params[j] = parseType(signature.Params().At(j).Type())
			}
		}

		returns := []Type{}
		if signature.Results() != nil {
			returns = make([]Type, signature.Results().Len())
			for j := 0; j < signature.Results().Len(); j++ {
				returns[j] = parseType(signature.Results().At(j).Type())
			}
		}

		_, isPointerReceiver := signature.Recv().Type().(*types.Pointer)

		methods[index] = Method{
			Name:              name,
			IsExported:        ast.IsExported(name),
			IsPointerReceiver: isPointerReceiver,
			Params:            params,
			Returns:           returns,
			Comments:          []string{}, // TODO
		}

		index++
	}

	// cleanup empty methods
	for i := 0; i < len(methods); i++ {
		if methods[i].Name == "" {
			methods = append(methods[:i], methods[i+1:]...)
			i--
		}
	}

	return methods, nil
}
