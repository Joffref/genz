package genz

import (
	"flag"
	"fmt"
	"github.com/Joffref/genz/internal/command"
	"github.com/Joffref/genz/internal/generator"
	"github.com/Joffref/genz/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type generateCommand struct {
}

const (
	generateCommandUsage = `Usage of genz:
	genz [flags] -type T -template foo.tmpl [directory]
	genz [flags] -type T -template foo.tmpl files... # Must be a single package
Flags:`
)

var (
	generateCmd  = flag.NewFlagSet("", flag.ExitOnError)
	typeNames    = generateCmd.String("type", "", "comma-separated list of type names; must be set")
	templateFile = generateCmd.String("template", "", "go-template file name")
	output       = generateCmd.String("output", "", "output file name; default srcdir/<type>.gen.go")
	buildTags    = generateCmd.String("tags", "", "comma-separated list of build tags to apply")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("genz: ")
	generateCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", generateCommandUsage)
		generateCmd.PrintDefaults()
	}
	command.RegisterCommand("generate", generateCommand{})
	command.SetRootCommand(generateCommand{})
}

func (c generateCommand) ValidateArgs() error {
	if len(*typeNames) == 0 {
		generateCmd.Usage()
		return fmt.Errorf("missing 'type' argument")
	}
	if len(*templateFile) == 0 {
		generateCmd.Usage()
		return fmt.Errorf("missing 'template' argument")
	}
	return nil
}

func (c generateCommand) FlagSet() *flag.FlagSet {
	return generateCmd
}

func (c generateCommand) Run() error {

	types := strings.Split(*typeNames, ",")

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	template, err := os.ReadFile(*templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %v", *templateFile, err)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	g := generator.Generator{
		Template: string(template),
	}

	g.ParsePackage(args, tags)
	// Run generate for each type.
	for _, typeName := range types {
		g.Generate(typeName)
	}

	src := g.Format()

	var dir string
	if len(args) == 1 && utils.IsDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			return fmt.Errorf("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s.gen.go", types[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}

	if err = os.WriteFile(outputName, src, 0644); err != nil {
		return fmt.Errorf("writing output: %s", err)
	}

	log.Printf("wrote %s (%d bytes)", outputName, len(src))
	return nil
}
