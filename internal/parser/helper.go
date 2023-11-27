package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
)

func parseElementType(pkg *packages.Package, name string) (Element, error) {
	if pkg.Types == nil {
		return Element{}, fmt.Errorf("package %s has no types", pkg.Name)
	}
	return Element{
		Type: Type{
			Name:         fmt.Sprintf("%s.%s", pkg.Name, name),
			InternalName: name,
		},
	}, nil
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

func parseType(t types.Type) Type {
	// Remove every qualifier before the type name
	// transforming "github.com/google/uuid.UUID" into "UUID"
	noPackageQualifier := func(_ *types.Package) string { return "" }

	// Adds the package name qualifier before the type name
	// transforming "github.com/google/uuid.UUID" into "uuid.UUID"
	packageNameQualifier := func(pkg *types.Package) string {
		return pkg.Name()
	}

	return Type{
		Name:         types.TypeString(t, packageNameQualifier), // (e.g. "uuid.UUID")
		InternalName: types.TypeString(t, noPackageQualifier),   // (e.g. "UUID")
	}

}

func loadAstExpr(pkg *packages.Package, typeName string) (ast.Expr, error) {
	for ident := range pkg.TypesInfo.Defs {
		if ident.Name == typeName {
			switch ident.Obj.Decl.(type) {
			case *ast.TypeSpec:
				return ident.Obj.Decl.(*ast.TypeSpec).Type, nil
			default:
				return nil, fmt.Errorf("%s is not a type", ident.Name)
			}
		}
	}
	return nil, fmt.Errorf("%s not found in package %s", typeName, pkg.Name)
}
