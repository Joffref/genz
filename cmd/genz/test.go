package genz

import (
	"flag"
	"fmt"
	"github.com/Joffref/genz/internal/command"
	"github.com/Joffref/genz/internal/testing"
	"log"
)

const (
	testUsage = "run tests"
)

type testCommand struct {
}

var (
	test      = flag.NewFlagSet("test", flag.ExitOnError)
	directory = test.String("directory", ".", "directory where to run the tests")
	exitOnErr = test.Bool("exit-on-error", false, "exit on first error")
	verbose   = test.Bool("verbose", false, "verbose mode")
)

func (t testCommand) FlagSet() *flag.FlagSet {
	return test
}

func (t testCommand) Run() error {
	if err := testing.RunTests(*directory, *verbose, *exitOnErr); err != nil {
		return fmt.Errorf("test(s) exited with error(s)")
	}
	return nil
}

func (t testCommand) ValidateArgs() error {
	return nil
}

func init() {
	test.BoolVar(verbose, "v", false, "verbose mode") // alias
	test.Usage = func() {
		log.Printf("%s\n", testUsage)
	}
	command.RegisterCommand("test", testCommand{})
}
