package parser

import (
	"github.com/Joffref/genz/pkg/models"
	"go/ast"
	"go/types"
)

// parseMethods returns a slice of models.Method from the given map of method signatures.
// It's not responsible for parsing the methods' comments. It should be done by the caller.
func parseMethods(signatures map[string]*types.Signature) ([]models.Method, error) {
	methods := make([]models.Method, len(signatures))

	index := 0
	for name, signature := range signatures {

		params := []models.Type{}
		if signature.Params() != nil {
			params = make([]models.Type, signature.Params().Len())
			for j := 0; j < signature.Params().Len(); j++ {
				params[j] = parseType(signature.Params().At(j).Type())
			}
		}

		returns := []models.Type{}
		if signature.Results() != nil {
			returns = make([]models.Type, signature.Results().Len())
			for j := 0; j < signature.Results().Len(); j++ {
				returns[j] = parseType(signature.Results().At(j).Type())
			}
		}

		_, isPointerReceiver := signature.Recv().Type().(*types.Pointer)

		methods[index] = models.Method{
			Name:              name,
			IsExported:        ast.IsExported(name),
			IsPointerReceiver: isPointerReceiver,
			Params:            params,
			Returns:           returns,
			Comments:          []string{}, // It has to be filled later by the caller.
		}

		index++
	}

	// cleanup empty methods
	// TODO: check if this is necessary
	for i := 0; i < len(methods); i++ {
		if methods[i].Name == "" {
			methods = append(methods[:i], methods[i+1:]...)
			i--
		}
	}

	return methods, nil
}
