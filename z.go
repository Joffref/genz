package main

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	ZType struct {
		Type Type

		Attributes []Attributes
	}
	Attributes struct {
		Name string
		Type Type

		Keys []map[string]string
	}
	Type string
)

func NewZType(typeName string, pkg *packages.Package) (ZType, error) {
	z := ZType{
		Type: Type(typeName),
	}

	for ident, obj := range pkg.TypesInfo.Defs {
		if ident.Name == typeName {
			structType, isStruct := obj.Type().Underlying().(*types.Struct)
			if !isStruct {
				return ZType{}, fmt.Errorf("type %s is not a struct", typeName)
			}

			z.Attributes = readStructAttributes(structType)
		}
	}

	return z, nil
}

func readStructAttributes(structType *types.Struct) []Attributes {
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
