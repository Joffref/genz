package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
)

func parseInterface(pkg *packages.Package, interfaceName string, interfaceType *ast.InterfaceType) (Element, error) {
	parsedInterface, err := parseElementType(pkg, interfaceName)
	if err != nil {
		return Element{}, err
	}
	signatures := map[string]*types.Signature{}
	comments := map[string][]string{}
	for _, method := range interfaceType.Methods.List {
		switch pkg.TypesInfo.TypeOf(method.Type).(type) {
		case *types.Signature:
			signatures[method.Names[0].Name] = pkg.TypesInfo.TypeOf(method.Type).(*types.Signature)
			if method.Doc != nil {
				for _, comment := range method.Doc.List {
					fmt.Println(comment.Text)
					comments[method.Names[0].Name] = append(comments[method.Names[0].Name], comment.Text[2:])
				}
			}
		case *types.Named: // Embedded interface
			namedType := pkg.TypesInfo.TypeOf(method.Type).(*types.Named)
			iface := namedType.Origin().Underlying().(*types.Interface).Complete()
			for i := 0; i < iface.NumMethods(); i++ {
				signatures[iface.Method(i).Name()] = iface.Method(i).Type().(*types.Signature)
				if method.Doc != nil {
					for _, comment := range method.Doc.List {
						comments[iface.Method(i).Name()] = append(comments[iface.Method(i).Name()], comment.Text[2:])
					}
				}
			}
		}

	}
	methods, err := parseMethods(signatures)
	if err != nil {
		return Element{}, err
	}
	var methodsWithComments []Method
	for _, method := range methods {
		if comments[method.Name] != nil {
			addCommentsToMethod(&method, comments[method.Name])
		}
		methodsWithComments = append(methodsWithComments, method)
	}
	parsedInterface.Methods = methodsWithComments
	return parsedInterface, nil
}

func addCommentsToMethod(method *Method, comments []string) {
	method.Comments = comments
}
