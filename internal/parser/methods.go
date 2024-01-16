package parser

import (
	"go/ast"
	"go/doc"
	"go/types"
	"strconv"
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

	params := map[string]models.Type{}
	if signature.Params() != nil {
		params = make(map[string]models.Type, signature.Params().Len())
		for j := 0; j < signature.Params().Len(); j++ {
			if signature.Params().At(j).Name() == "" { // unnamed parameter
				params[strconv.Itoa(j)] = parseType(signature.Params().At(j).Type())
				continue
			}
			params[signature.Params().At(j).Name()] = parseType(signature.Params().At(j).Type())
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
