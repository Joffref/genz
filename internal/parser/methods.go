package parser

import (
	"go/ast"
	"go/doc"
	"go/types"
	"strings"

	"github.com/Joffref/genz/pkg/models"
)

func parseMethod(name string, signature *types.Signature) (models.Method, error) {
	return parseMethodWithComments(nil, name, signature)
}

func parseMethodWithComments(doc *doc.Func, name string, signature *types.Signature) (models.Method, error) {
	comments := []string{}
	if doc != nil && doc.Doc != "" {
		comments = strings.Split(strings.Trim(doc.Doc, "\n"), "\n")
	}

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
		Comments:          comments,
	}, nil
}
