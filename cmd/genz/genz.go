package genz

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/pkg/models"
	"io"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/Joffref/genz/internal/command"
	"github.com/Joffref/genz/internal/utils"
)

type generateCommand struct {
}

const (
	generateCommandUsage = `Usage of genz:
	genz [flags] -type T -generator path/to/my/custom/generator:entrypoint [directory]
	genz [flags] -type T -generator path/to/my/custom/generator:entrypoint files... # Must be a single package
Flags:`
)

var (
	generateCmd = flag.NewFlagSet("", flag.ExitOnError)
	typeName    = generateCmd.String("type", "", "name of the type to parse")
	generator   = generateCmd.String("generator", "", "name of the generator to use (e.g: path/to/my/custom/generator:entrypoint)")
	output      = generateCmd.String("output", "", "output file name; default srcdir/<type>.gen.go")
	buildTags   = generateCmd.String("tags", "", "comma-separated list of build tags to apply")
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
	if len(*typeName) == 0 {
		generateCmd.Usage()
		return fmt.Errorf("missing 'type' argument")
	}
	if len(*generator) == 0 {
		generateCmd.Usage()
		return fmt.Errorf("missing 'generator' argument")
	}
	if len(strings.Split(*generator, ":")) != 2 {
		generateCmd.Usage()
		return fmt.Errorf("invalid 'generator' argument: should be path/to/my/custom/generator:entrypoint")
	}
	return nil
}

func (c generateCommand) FlagSet() *flag.FlagSet {
	return generateCmd
}

func (c generateCommand) Run() error {
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	pluginPath := strings.Split(*generator, ":")[0]
	if _, err := os.Stat(pluginPath); err != nil {
		return fmt.Errorf("could not find generator plugin: %s", err)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("could not open generator plugin: %s", err)
	}

	userDefinedGenerator, err := p.Lookup(strings.Split(*generator, ":")[1])
	if err != nil {
		return fmt.Errorf("could not find generator entrypoint: %s", err)
	}

	userDefinedGeneratorFunc, ok := userDefinedGenerator.(func(writer io.Writer, element models.ParsedElement) error)
	if !ok {
		return fmt.Errorf("could not find generator entrypoint: %s", err)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := generateCmd.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	parsedElement, err := parser.Parser(utils.LoadPackage(args, tags), *typeName)

	var buf bytes.Buffer
	err = userDefinedGeneratorFunc(&buf, parsedElement)
	if err != nil {
		return fmt.Errorf("generating code: %s", err)
	}
	src := utils.Format(buf)

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
		baseName := fmt.Sprintf("%s.gen.go", *typeName)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}

	if err := os.WriteFile(outputName, src, 0644); err != nil {
		return fmt.Errorf("writing output: %s", err)
	}

	log.Printf("wrote %s (%d bytes)", outputName, len(src))
	return nil
}
