//nolint:unused
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
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
	template string
	buf      bytes.Buffer // Accumulated output.
	pkg      *Package     // Package we are scanning.
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
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

func (g *Generator) generate(typeName string) {
	log.Printf("generating template for type %s", typeName)

	tmpl, err := template.New("template").Parse(g.template)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	err = tmpl.Execute(&g.buf, nil)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	log.Printf("generated buffer (%d bytes)", g.buf.Len())
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
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

var (
	typeNames    = flag.String("type", "", "comma-separated list of type names; must be set")
	templateFile = flag.String("template", "", "go-template file name")
	output       = flag.String("output", "", "output file name; default srcdir/<type>.gen.go")
	buildTags    = flag.String("tags", "", "comma-separated list of build tags to apply")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of genz:\n")
	fmt.Fprintf(os.Stderr, "\tgenz [flags] -type T -template foo.tmpl [directory]\n")
	fmt.Fprintf(os.Stderr, "\tgenz [flags] -type T -template foo.tmpl files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("genz: ")
	flag.Usage = Usage
	flag.Parse()

	if len(*typeNames) == 0 {
		flag.Usage()
		log.Fatal("missing 'type' argument")
	}
	types := strings.Split(*typeNames, ",")

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	if len(*templateFile) == 0 {
		flag.Usage()
		log.Fatal("missing 'template' argument")
	}
	template, err := os.ReadFile(*templateFile)
	if err != nil {
		log.Fatalf("failed to read template file %s: %v", *templateFile, err)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	g := Generator{
		template: string(template),
	}

	g.parsePackage(args, tags)
	// Run generate for each type.
	for _, typeName := range types {
		g.generate(typeName)
	}

	src := g.format()

	var dir string
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s.gen.go", types[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}

	err = os.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}

	log.Printf("wrote %s (%d bytes)", outputName, len(src))
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
