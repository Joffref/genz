package parser

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/types"
	"strings"

	"github.com/Joffref/genz/pkg/models"

	"golang.org/x/tools/go/packages"
)

func parseStruct(pkg *packages.Package, structName string, structType *ast.StructType) (models.Element, error) {
	parsedStruct := models.Element{
		Type: models.Type{
			Name:         fmt.Sprintf("%s.%s", pkg.Name, structName),
			InternalName: structName,
		},
	}

	attributes, err := structAttributes(pkg.TypesInfo, structType)
	if err != nil {
		return models.Element{}, err
	}

	parsedStruct.Attributes = attributes

	pkgDoc, err := doc.NewFromFiles(pkg.Fset, pkg.Syntax, "./", doc.AllDecls)
	if err != nil {
		return models.Element{}, fmt.Errorf("failed to create doc while parsing struct: %w", err)
	}

	var methods []models.Method

	for ident, object := range pkg.TypesInfo.Uses {
		if ident.Name == structName {
			namedType, err := objectAsNamedType(object)
			if err != nil {
				return models.Element{}, err
			}

			for i := 0; i < namedType.NumMethods(); i++ {
				method, err := parseMethodWithComments(
					getFuncDoc(pkgDoc, structName, namedType.Method(i).Name()),
					namedType.Method(i).Name(),
					namedType.Method(i).Type().(*types.Signature),
				)
				if err != nil {
					return models.Element{}, err
				}
				methods = append(methods, method)
			}
		}
	}
	parsedStruct.Methods = methods
	return parsedStruct, nil
}

func structAttributes(typesInfo *types.Info, structType *ast.StructType) ([]models.Attribute, error) {
	attributes := make([]models.Attribute, len(structType.Fields.List))

	for i, field := range structType.Fields.List {
		comments := []string{}
		if field.Doc != nil {
			for _, comment := range field.Doc.List {
				comments = append(comments, comment.Text[2:])
			}
		}
		attributes[i] = models.Attribute{
			Name:     field.Names[0].Name,
			Type:     parseType(typesInfo.TypeOf(field.Type)),
			Comments: comments,
		}
		if field.Tag != nil {
			tags, err := parseTags(field.Tag.Value)
			if err != nil {
				return nil, err
			}
			attributes[i].Tags = tags
		}
	}

	return attributes, nil
}

// parseTags take a string of tags (e.g. `json:"name,omitempty" xml:"name"`)
// and returns a map of tags (e.g. map[string]string{"json": "name,omitempty", "xml": "name"})
func parseTags(tags string) (map[string]string, error) {
	var result = make(map[string]string)
	tags = strings.ReplaceAll(tags, "`", "")
	for _, tag := range strings.Split(tags, "\" ") {
		if tag == "" {
			continue
		}
		splitTag := strings.Split(tag, ":")
		if len(splitTag) != 2 {
			return nil, fmt.Errorf("invalid tag: %s", tag)
		}
		result[splitTag[0]] = strings.ReplaceAll(splitTag[1], "\"", "")
	}
	return result, nil
}

func getFuncDoc(pkgDoc *doc.Package, structName, funcName string) *doc.Func {
	var typeDoc *doc.Type
	for _, t := range pkgDoc.Types {
		if t.Name == structName {
			typeDoc = t
		}
	}

	for _, m := range typeDoc.Methods {
		if m.Name == funcName {
			return m
		}
	}

	return nil
}
