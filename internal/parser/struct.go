package parser

import (
	"fmt"
	"github.com/Joffref/genz/pkg/models"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

func parseStruct(pkg *packages.Package, structName string, structType *ast.StructType) (models.Element, error) {
	parsedStruct := models.Element{
		Type: models.Type{
			Name:         fmt.Sprintf("%s.%s", pkg.Name, structName),
			InternalName: structName,
		},
	}

	parsedStruct.Attributes = structAttributes(pkg.TypesInfo, structType)

	var methods []models.Method

	for ident, object := range pkg.TypesInfo.Uses {
		if ident.Name == structName {
			namedType, err := objectAsNamedType(object)
			if err != nil {
				return models.Element{}, err
			}

			for i := 0; i < namedType.NumMethods(); i++ {
				method, err := parseMethod(namedType.Method(i).Name(), namedType.Method(i).Type().(*types.Signature))
				if err != nil {
					return models.Element{}, err
				}
				methods = append(methods, method)
			}
		}
	}
	parsedStruct.Methods = methods
	return parsedStruct, nil
}

func structAttributes(typesInfo *types.Info, structType *ast.StructType) []models.Attribute {
	attributes := make([]models.Attribute, len(structType.Fields.List))

	for i, field := range structType.Fields.List {
		comments := []string{}
		if field.Doc != nil {
			for _, comment := range field.Doc.List {
				comments = append(comments, comment.Text[2:])
			}
		}
		attributes[i] = models.Attribute{
			Name:     field.Names[0].Name,
			Type:     parseType(typesInfo.TypeOf(field.Type)),
			Comments: comments,
		}
		if field.Tag != nil {
			attributes[i].Tags = parseTags(field.Tag.Value)
		}
	}

	return attributes
}

// parseTags take a string of tags (e.g. `json:"name,omitempty" xml:"name"`)
// and returns a map of tags (e.g. map[string]string{"json": "name,omitempty", "xml": "name"})
func parseTags(tags string) map[string]string {
	var result = make(map[string]string)
	tags = strings.ReplaceAll(tags, "`", "")
	for _, tag := range strings.Split(tags, "\" ") {
		if tag == "" {
			continue
		}
		splitTag := strings.Split(tag, ":")
		result[splitTag[0]] = strings.ReplaceAll(splitTag[1], "\"", "")
	}
	return result
}
