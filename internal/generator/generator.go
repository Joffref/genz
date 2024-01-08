//nolint:unused
package generator

import (
	"bytes"
	"fmt"
	"github.com/Joffref/genz/pkg/models"
	"html/template"
	"log"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/tools/go/packages"
)

type parseFunc func(pkg *packages.Package, typeName string) (models.ParsedElement, error)

func Generate(
	pkg *packages.Package,
	templateContent string,
	typeName string,
	parse parseFunc,
) (bytes.Buffer, error) {
	log.Printf("generating template for type %s", typeName)

	parsedElement, err := parse(pkg, typeName)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to inspect package: %v", err)
	}

	tmpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(templateContent)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to parse template: %v", err)
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, parsedElement)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute template: %v", err)
	}

	log.Printf("generated buffer (%d bytes)", buf.Len())
	return buf, nil
}
