package parser

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

func Parser(pkg *packages.Package, typeName string) (ParsedElement, error) {
	parsedElement, err := parsePackage(pkg)
	if err != nil {
		return ParsedElement{}, err
	}
	expr, err := loadAstExpr(pkg, typeName)
	if err != nil {
		return ParsedElement{}, err
	}
	element, err := parseElement(pkg, expr, typeName)
	if err != nil {
		return ParsedElement{}, err
	}
	parsedElement.Element = element
	return parsedElement, nil
}

func parseElement(pkg *packages.Package, expr ast.Expr, name string) (Element, error) {
	switch expr := expr.(type) {
	case *ast.StructType:
		return parseStruct(pkg, name, expr)
	case *ast.InterfaceType:
		return parseInterface(pkg, name, expr)
	default:
		return Element{}, fmt.Errorf("unsupported type %T", expr)
	}
}
