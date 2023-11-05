//nolint:unused
package generator

import (
	"bytes"
	"go/format"
	"log"
	"text/template"

	"github.com/Joffref/genz/internal/parser"
	"golang.org/x/tools/go/packages"
)

type Generator struct {
	Template string
	Pkg      *packages.Package // Package we are scanning.
	buf      bytes.Buffer      // Accumulated output.
}

func (g *Generator) Generate(typeName string) {
	log.Printf("generating template for type %s", typeName)

	parsedType, err := parser.Parse(g.Pkg, typeName)
	if err != nil {
		log.Fatalf("failed to inspect package: %v", err)
	}

	tmpl, err := template.New("template").Parse(g.Template)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	err = tmpl.Execute(&g.buf, parsedType)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	log.Printf("generated buffer (%d bytes)", g.buf.Len())
}

func (g *Generator) Format() []byte {
	log.Print("gofmt-ing buffer")

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
