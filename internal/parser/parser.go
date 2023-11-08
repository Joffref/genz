package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func Parse(pkg *packages.Package, structName string) (Struct, error) {
	for ident := range pkg.TypesInfo.Defs {
		if ident.Name == structName {
			structType, err := identAsStructType(ident)
			if err != nil {
				return Struct{}, err
			}

			return Struct{
				Type:       Type(structName),
				Attributes: structAttributes(pkg.TypesInfo, structType),
			}, nil
		}
	}
	return Struct{}, fmt.Errorf("struct %s not found in package %s", structName, pkg.Name)
}

func identAsStructType(ident *ast.Ident) (*ast.StructType, error) {
	typeSpec, isTypeSpec := ident.Obj.Decl.(*ast.TypeSpec)
	if !isTypeSpec {
		return nil, fmt.Errorf("%s is not a type", ident.Name)
	}

	structDeclaration, isStruct := typeSpec.Type.(*ast.StructType)
	if !isStruct {
		return nil, fmt.Errorf("%s is not a struct", ident.Name)
	}

	return structDeclaration, nil
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
			Type:     Type(typesInfo.TypeOf(field.Type).String()),
			Comments: comments,
		}
	}

	return attributes
}
