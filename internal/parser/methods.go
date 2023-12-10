package parser

import (
	"github.com/Joffref/genz/pkg/models"
	"go/ast"
	"go/types"
)

// parseMethod returns a  models.Method from the given map of method signatures.
// It's not responsible for parsing the methods' comments. It should be done by the caller.
func parseMethod(name string, signature *types.Signature) (models.Method, error) {

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

	return models.Method{
		Name:              name,
		IsExported:        ast.IsExported(name),
		IsPointerReceiver: isPointerReceiver,
		Params:            params,
		Returns:           returns,
		Comments:          []string{}, // It has to be filled later by the caller.
	}, nil
}
