package parser

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func Parse(pkg *packages.Package, typeName string) (ParsedType, error) {
	p := ParsedType{
		Type: Type(typeName),
	}

	for ident := range pkg.TypesInfo.Defs {
		if ident.Name == typeName {
			typeSpec, isTypeSpec := ident.Obj.Decl.(*ast.TypeSpec)
			if !isTypeSpec {
				return ParsedType{}, fmt.Errorf("%s is not a type", typeName)
			}
			structDeclaration, isStruct := typeSpec.Type.(*ast.StructType)
			if !isStruct {
				return ParsedType{}, fmt.Errorf("%s is not a struct", typeName)
			}

			for _, field := range structDeclaration.Fields.List {
				comments := make([]string, len(field.Doc.List))
				for i, comment := range field.Doc.List {
					comments[i] = comment.Text[2:]
				}
				p.Attributes = append(p.Attributes, Attributes{
					Name:     field.Names[0].Name,
					Type:     Type(pkg.TypesInfo.TypeOf(field.Type).String()),
					Comments: comments,
				})
			}
		}
	}

	return p, nil
}
