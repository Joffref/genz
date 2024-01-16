package genz

import (
	"flag"
	"fmt"
	"github.com/Joffref/genz/internal/parser"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Joffref/genz/internal/command"
	"github.com/Joffref/genz/internal/generator"
	"github.com/Joffref/genz/internal/utils"
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
	generateCmd      = flag.NewFlagSet("", flag.ExitOnError)
	typeName         = generateCmd.String("type", "", "name of the type to parse")
	templateLocation = generateCmd.String("template", "", "go-template local or remote file")
	output           = generateCmd.String("output", "", "output file name; default srcdir/<type>.gen.go")
	buildTags        = generateCmd.String("tags", "", "comma-separated list of build tags to apply")
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
	if len(*templateLocation) == 0 {
		generateCmd.Usage()
		return fmt.Errorf("missing 'template' argument")
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

	var template []byte
	if url, _ := url.ParseRequestURI(*templateLocation); url != nil {
		response, err := http.Get(*templateLocation)
		if err != nil {
			return fmt.Errorf("failed to make a request to %s: %v", *templateLocation, err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("could not read body of remote template %s: %v", *templateLocation, err)
		}
		template = body
	} else {
		file, err := os.ReadFile(*templateLocation)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %v", *templateLocation, err)
		}
		template = file
	}

	// We accept either one directory or a list of files. Which do we have?
	args := generateCmd.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}
	buf, err := generator.Generate(
		utils.LoadPackage(args, tags),
		string(template),
		*typeName,
		parser.Parse,
	)
	if err != nil {
		return err
	}

	src := generator.Format(buf)

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
