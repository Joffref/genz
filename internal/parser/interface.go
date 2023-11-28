package parser

import (
	"github.com/Joffref/genz/pkg/models"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
)

// parseInterface parses the given interface and returns a models.Element.
// It also parses the subInterfaces of the given interface and the methods of the subInterfaces recursively.
func parseInterface(pkg *packages.Package, interfaceName string, interfaceType *ast.InterfaceType) (models.Element, error) {

	parsedInterface, err := parseElementType(pkg, interfaceName)
	if err != nil {
		return models.Element{}, err
	}

	signatures := map[string]*types.Signature{}
	comments := map[string][]string{}

	for _, method := range interfaceType.Methods.List {
		switch pkg.TypesInfo.TypeOf(method.Type).(type) {
		case *types.Signature:
			signatures[method.Names[0].Name] = pkg.TypesInfo.TypeOf(method.Type).(*types.Signature)
			if method.Doc != nil {
				for _, comment := range method.Doc.List {
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
		return models.Element{}, err
	}

	var methodsWithComments []models.Method
	for _, method := range methods {
		if comments[method.Name] != nil {
			addCommentsToMethod(&method, comments[method.Name])
		}
		if len(methodsWithComments) == 0 {
			methodsWithComments = append(methodsWithComments, method)
			continue
		}
		// order methods by name, to make the output deterministic.
		for i, methodWithComments := range methodsWithComments {
			if method.Name < methodWithComments.Name { // It's not possible to have two methods with the same name in an interface.
				methodsWithComments = append(methodsWithComments[:i], append([]models.Method{method}, methodsWithComments[i:]...)...)
				break
			}
			if i == len(methodsWithComments)-1 {
				methodsWithComments = append(methodsWithComments, method)
				break
			}
		}
	}
	parsedInterface.Methods = methodsWithComments
	return parsedInterface, nil
}
