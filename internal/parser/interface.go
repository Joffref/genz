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

	var methods []models.Method

	for _, method := range interfaceType.Methods.List {
		switch pkg.TypesInfo.TypeOf(method.Type).(type) {
		case *types.Signature:
			methodModel, err := parseMethod(method.Names[0].Name, pkg.TypesInfo.TypeOf(method.Type).(*types.Signature))
			if err != nil {
				return models.Element{}, err
			}
			if method.Doc != nil {
				for _, comment := range method.Doc.List {
					methodModel.Comments = append(methodModel.Comments, comment.Text[2:])
				}
			}
			methods = append(methods, methodModel)
		case *types.Named: // Embedded interface
			namedType := pkg.TypesInfo.TypeOf(method.Type).(*types.Named)
			iface := namedType.Origin().Underlying().(*types.Interface).Complete()
			for i := 0; i < iface.NumMethods(); i++ {
				methodModel, err := parseMethod(iface.Method(i).Name(), iface.Method(i).Type().(*types.Signature))
				if err != nil {
					return models.Element{}, err
				}
				if method.Doc != nil {
					for _, comment := range method.Doc.List {
						methodModel.Comments = append(methodModel.Comments, comment.Text[2:])
					}
				}
				methods = append(methods, methodModel)
			}
		}

	}
	parsedInterface.Methods = methods
	return parsedInterface, nil
}
