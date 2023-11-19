//nolint:unused
package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"

	"github.com/Joffref/genz/internal/parser"
	"github.com/Masterminds/sprig/v3"
	"golang.org/x/tools/go/packages"
)

type parseFunc func(pkg *packages.Package, structName string) (parser.Struct, error)

func Generate(
	pkg *packages.Package,
	templateContent string,
	typeName string,
	parse parseFunc,
) (bytes.Buffer, error) {
	log.Printf("generating template for type %s", typeName)

	parsedType, err := parse(pkg, typeName)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to inspect package: %v", err)
	}

	tmpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(templateContent)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to parse template: %v", err)
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, parsedType)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute template: %v", err)
	}

	log.Printf("generated buffer (%d bytes)", buf.Len())
	return buf, nil
}

func Format(buf bytes.Buffer) []byte {
	log.Print("gofmt-ing buffer")

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return buf.Bytes()
	}
	return src
}
