package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func parseStruct(pkg *packages.Package, structName string, structType *ast.StructType) (Element, error) {
	parsedStruct := Element{
		Type: Type{
			Name:         fmt.Sprintf("%s.%s", pkg.Name, structName),
			InternalName: structName,
		},
	}

	parsedStruct.Attributes = structAttributes(pkg.TypesInfo, structType)

	for ident, object := range pkg.TypesInfo.Uses {
		if ident.Name == structName {
			namedType, err := objectAsNamedType(object)
			if err != nil {
				return Element{}, err
			}

			signatures := map[string]*types.Signature{}
			for i := 0; i < namedType.NumMethods(); i++ {
				signatures[namedType.Method(i).Name()] = namedType.Method(i).Type().(*types.Signature)
			}

			parsedStruct.Methods, err = parseMethods(signatures)
			if err != nil {
				return Element{}, err
			}
		}
	}

	return parsedStruct, nil
}

func structAttributes(typesInfo *types.Info, structType *ast.StructType) []Attribute {
	attributes := make([]Attribute, len(structType.Fields.List))

	for i, field := range structType.Fields.List {
		comments := []string{}
		if field.Doc != nil {
			for _, comment := range field.Doc.List {
				comments = append(comments, comment.Text[2:])
			}
		}

		attributes[i] = Attribute{
			Name:     field.Names[0].Name,
			Type:     parseType(typesInfo.TypeOf(field.Type)),
			Comments: comments,
		}
	}

	return attributes
}
