package cli

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Joffref/genz/genz"
	"os"
)

type Command interface {
	FlagSet() *flag.FlagSet
	Execute() error
	GetConfig() ([]byte, error)
	ValidateArgs() error
}

var (
	defaultCmd = flag.NewFlagSet("", flag.ExitOnError)
	typeName   = defaultCmd.String("type", "", "name of the type to parse")
	output     = defaultCmd.String("output", "", "output file name; default srcdir/<type>.gen.go")
	config     = defaultCmd.String("config", "", "path to the config file, if any")
)

type commandFromGenerator struct {
	generator genz.Generator
	name      string
}

func NewCommandFromGenerator(name string, generator genz.Generator) Command {
	parseFlags()
	return commandFromGenerator{
		generator: generator,
		name:      name,
	}
}

func (c commandFromGenerator) FlagSet() *flag.FlagSet {
	return defaultCmd
}

func (c commandFromGenerator) GetConfig() ([]byte, error) {
	file, err := os.ReadFile(*config)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c commandFromGenerator) Execute() error {
	if err := c.ValidateArgs(); err != nil {
		return err
	}
	if len(*output) == 0 {
		*output = fmt.Sprintf("%s.gen.go", *typeName)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := defaultCmd.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	parsedElement, err := genz.Parse(*typeName, args, *output)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = c.generator(&buf, parsedElement); err != nil {
		return err
	}

	file, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err = file.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (c commandFromGenerator) ValidateArgs() error {
	if len(*typeName) == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", c.name)
		defaultCmd.PrintDefaults()
		return fmt.Errorf("missing 'type' argument")
	}
	return nil
}

func parseFlags() {
	if err := defaultCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
