package parser

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func Parse(pkg *packages.Package, structName string) (Struct, error) {
	parsedStruct := Struct{Type: Type(structName)}
	found := false
	for ident := range pkg.TypesInfo.Defs {
		if ident.Name == structName {
			structType, err := identAsStructType(ident)
			if err != nil {
				return Struct{}, err
			}

			parsedStruct.Attributes = structAttributes(pkg.TypesInfo, structType)
			found = true
			break
		}
	}
	if !found {
		return Struct{}, fmt.Errorf("struct %s not found in package %s", structName, pkg.Name)
	}

	for ident, object := range pkg.TypesInfo.Uses {
		if ident.Name == structName {
			namedType, err := objectAsNamedType(object)
			if err != nil {
				return Struct{}, err
			}

			parsedStruct.Methods, err = structMethods(namedType)
			if err != nil {
				return Struct{}, err
			}
		}
	}

	return parsedStruct, nil
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

func objectAsNamedType(object types.Object) (*types.Named, error) {
	typeName, isTypeName := object.(*types.TypeName)
	if !isTypeName {
		return nil, fmt.Errorf("%s is not a TypeName", object.Name())
	}
	namedType, isNamedType := typeName.Type().(*types.Named)
	if !isNamedType {
		return nil, fmt.Errorf("%s is not a named type", object.Name())
	}

	return namedType, nil
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

func structMethods(namedType *types.Named) ([]Method, error) {
	methods := make([]Method, namedType.NumMethods())

	for i := 0; i < namedType.NumMethods(); i++ {
		declaration := namedType.Method(i)
		signature, isSignature := declaration.Type().(*types.Signature)
		if !isSignature {
			return nil, fmt.Errorf("cannot get signature from method %s", declaration.Name())
		}

		params := []Type{}
		if signature.Params() != nil {
			params = make([]Type, signature.Params().Len())
			for j := 0; j < signature.Params().Len(); j++ {
				params[j] = Type(signature.Params().At(j).Type().String())
			}
		}

		returns := []Type{}
		if signature.Results() != nil {
			returns = make([]Type, signature.Results().Len())
			for j := 0; j < signature.Results().Len(); j++ {
				returns[j] = Type(signature.Results().At(j).Type().String())
			}
		}

		_, isPointerReceiver := signature.Recv().Type().(*types.Pointer)

		methods[i] = Method{
			Name:              declaration.Name(),
			IsExported:        declaration.Exported(),
			IsPointerReceiver: isPointerReceiver,
			Params:            params,
			Returns:           returns,
			Comments:          []string{}, // TODO
		}
	}

	return methods, nil
}
