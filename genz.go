//nolint:unused
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type Generator struct {
	template string
	buf      bytes.Buffer // Accumulated output.
	pkg      *packages.Package
}

func (g *Generator) generate(typeName string) {
	log.Printf("generating template for type %s", typeName)

	zType, err := NewZType(typeName, g.pkg)
	if err != nil {
		log.Fatalf("failed to inspect package: %v", err)
	}

	tmpl, err := template.New("template").Parse(g.template)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	err = tmpl.Execute(&g.buf, zType)
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

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func parsePackage(patterns []string, tags []string) *packages.Package {
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

	return pkgs[0]
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
		log.Fatalf("failed to read template file: %s", *templateFile)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	g := Generator{
		template: string(template),
		pkg:      parsePackage(args, tags),
	}

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
