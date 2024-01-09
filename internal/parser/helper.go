package parser

import (
	"fmt"
	"github.com/Joffref/genz/genz/models"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
)

// parseElementType initializes a models.Element with the given name and package.
// It does not parse the element's methods or attributes.
// It should be called first by every parse* function.
func parseElementType(pkg *packages.Package, name string) (models.Element, error) {
	if pkg.Types == nil {
		return models.Element{}, fmt.Errorf("package %s has no types", pkg.Name)
	}
	return models.Element{
		Type: models.Type{
			Name:         fmt.Sprintf("%s.%s", pkg.Name, name),
			InternalName: name,
		},
	}, nil
}

// objectAsNamedType returns the given object as a *types.Named.
// It returns an error if the object is not a *types.TypeName or if the type of the *types.TypeName is not a *types.Named.
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

// parseType returns a models.Type from the given types.Type.
// It returns the type name with the package qualifier and without the package qualifier.
func parseType(t types.Type) models.Type {
	// Remove every qualifier before the type name
	// transforming "github.com/google/uuid.UUID" into "UUID"
	noPackageQualifier := func(_ *types.Package) string { return "" }

	// Adds the package name qualifier before the type name
	// transforming "github.com/google/uuid.UUID" into "uuid.UUID"
	packageNameQualifier := func(pkg *types.Package) string {
		return pkg.Name()
	}

	return models.Type{
		Name:         types.TypeString(t, packageNameQualifier), // (e.g. "uuid.UUID")
		InternalName: types.TypeString(t, noPackageQualifier),   // (e.g. "UUID")
	}

}

// loadAstExpr returns the ast.Expr of the given typeName in the given package.
// It returns an error if the typeName is not found in the package.
// Be aware that ast.Expr is an interface, so the returned value can be of any type.
func loadAstExpr(pkg *packages.Package, typeName string) (ast.Expr, error) {
	for ident := range pkg.TypesInfo.Defs {
		if ident.Name == typeName {
			if ident.Obj == nil { // could be a name overlapping with a type.
				continue
			}
			switch ident.Obj.Decl.(type) {
			case *ast.TypeSpec:
				return ident.Obj.Decl.(*ast.TypeSpec).Type, nil
			default: // could be a name overlapping with a type.
				continue
			}
		}
	}
	return nil, fmt.Errorf("%s not found in package %s", typeName, pkg.Name)
}
