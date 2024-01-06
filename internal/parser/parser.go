package parser

import (
	"fmt"
	"go/ast"

	"github.com/Joffref/genz/pkg/models"
	"golang.org/x/tools/go/packages"
)

// Parser returns a models.ParsedElement from the given *packages.Package and type name.
// It's the main entry point for the different parsers (struct, interface, etc).
func Parser(pkg *packages.Package, typeName string) (models.ParsedElement, error) {
	parsedElement, err := parsePackage(pkg)
	if err != nil {
		return models.ParsedElement{}, err
	}
	expr, err := loadAstExpr(pkg, typeName)
	if err != nil {
		return models.ParsedElement{}, err
	}
	element, err := parseElement(pkg, expr, typeName)
	if err != nil {
		return models.ParsedElement{}, err
	}
	parsedElement.Element = element
	return parsedElement, nil
}

// parseElement parses the given ast.Expr and returns a models.Element.
// It calls the appropriate parse* function depending on the type of the given ast.Expr.
func parseElement(pkg *packages.Package, expr ast.Expr, name string) (models.Element, error) {
	switch expr := expr.(type) {
	case *ast.StructType:
		return parseStruct(pkg, name, expr)
	case *ast.InterfaceType:
		return parseInterface(pkg, name, expr)
	default:
		return models.Element{}, fmt.Errorf("unsupported type %T", expr)
	}
}
