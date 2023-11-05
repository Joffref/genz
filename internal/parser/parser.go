package parser

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func Parse(pkg *packages.Package, typeName string) (ParsedType, error) {
	p := ParsedType{
		Type: Type(typeName),
	}

	for ident, obj := range pkg.TypesInfo.Defs {
		if ident.Name == typeName {
			structType, isStruct := obj.Type().Underlying().(*types.Struct)
			if !isStruct {
				return ParsedType{}, fmt.Errorf("type %s is not a struct", typeName)
			}

			p.Attributes = structAttributes(structType)
		}
	}

	return p, nil
}

func structAttributes(structType *types.Struct) []Attributes {
	var attributes []Attributes
	for i := 0; i < structType.NumFields(); i++ {

		attribute := structType.Field(i)
		attributes = append(attributes, Attributes{
			Name: attribute.Name(),
			Type: Type(attribute.Origin().Type().String()),
		})
	}

	return attributes
}
