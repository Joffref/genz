package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
)

func ParseInterface(pkg *packages.Package, interfaceName string) (interface{}, error) {
	parsedInterface := Interface{
		Type: Type{
			Name:         fmt.Sprintf("%s.%s", pkg.Name, interfaceName),
			InternalName: interfaceName,
		},
	}

	found := false
	for ident := range pkg.TypesInfo.Defs {
		if ident.Name == interfaceName {
			interfaceType, err := identAsInterfaceType(ident)
			if err != nil {
				return Interface{}, err
			}

			parsedInterface.Methods = interfaceMethods(pkg.TypesInfo, interfaceType)
			found = true
			break
		}
	}
	if !found {
		return Interface{}, fmt.Errorf("interface %s not found in package %s", interfaceName, pkg.Name)
	}

	return parsedInterface, nil
}

func identAsInterfaceType(ident *ast.Ident) (*ast.InterfaceType, error) {
	typeSpec, isTypeSpec := ident.Obj.Decl.(*ast.TypeSpec)
	if !isTypeSpec {
		return nil, fmt.Errorf("%s is not a type", ident.Name)
	}

	interfaceDeclaration, isInterface := typeSpec.Type.(*ast.InterfaceType)
	if !isInterface {
		return nil, fmt.Errorf("%s is not an interface", ident.Name)
	}

	return interfaceDeclaration, nil
}

func interfaceMethods(typesInfo *types.Info, interfaceType *ast.InterfaceType) []Method {
	methods := make([]Method, len(interfaceType.Methods.List))

	for i, method := range interfaceType.Methods.List {
		comments := []string{}
		if method.Doc != nil {
			for _, comment := range method.Doc.List {
				comments = append(comments, comment.Text[2:])
			}
		}

		signature, isSignature := typesInfo.TypeOf(method.Type).(*types.Signature)
		if !isSignature {
			panic(fmt.Errorf("cannot get signature from method %s", method.Names[0].Name))
		}

		params := []Type{}
		if signature.Params() != nil {
			params = make([]Type, signature.Params().Len())
			for j := 0; j < signature.Params().Len(); j++ {
				params[j] = parseType(signature.Params().At(j).Type())
			}
		}

		returns := []Type{}
		if signature.Results() != nil {
			returns = make([]Type, signature.Results().Len())
			for j := 0; j < signature.Results().Len(); j++ {
				returns[j] = parseType(signature.Results().At(j).Type())
			}
		}

		methods[i] = Method{
			Name:       method.Names[0].Name,
			IsExported: method.Names[0].IsExported(),
			Params:     params,
			Returns:    returns,
			Comments:   comments,
		}
	}

	return methods
}
