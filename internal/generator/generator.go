//nolint:unused
package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
	"text/template"
)

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName string
}

type Package struct {
	name      string
	pkgPath   string
	typesInfo *types.Info
	files     []*File
}

type Generator struct {
	Template string
	buf      bytes.Buffer // Accumulated output.
	pkg      *Package     // Package we are scanning.
}

// ParsePackage analyzes the single package constructed from the patterns and tags.
// ParsePackage exits if there is an error.
func (g *Generator) ParsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages matching %v", len(pkgs), strings.Join(patterns, " "))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	log.Printf("found package %s\n", pkg)
	g.pkg = &Package{
		name:      pkg.Name,
		pkgPath:   pkg.PkgPath,
		typesInfo: pkg.TypesInfo,
		files:     make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
}

func (g *Generator) Generate(typeName string) {
	log.Printf("generating template for type %s", typeName)

	tmpl, err := template.New("template").Parse(g.Template)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	err = tmpl.Execute(&g.buf, nil)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	log.Printf("generated buffer (%d bytes)", g.buf.Len())
}

// Format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) Format() []byte {
	log.Print("gofmt-ing buffer")

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
