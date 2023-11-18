package genz

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Joffref/genz/internal/command"
	"log"
	"os"
	"os/exec"
	"strings"
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

func init() {
	test.BoolVar(verbose, "v", false, "verbose mode") // alias
	test.Usage = func() {
		log.Printf("%s\n", testUsage)
	}
	command.RegisterCommand("test", testCommand{})
}

func (t testCommand) FlagSet() *flag.FlagSet {
	return test
}

func (t testCommand) Run() error {
	log.Printf("genz test %s\n", *directory)
	return runTests(*directory)
}

func (t testCommand) ValidateArgs() error {
	return nil
}

func runTests(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", directory, err)
	}
	var subDirectories []string
	for _, f := range entries {
		if f.IsDir() {
			subDirectories = append(subDirectories, f.Name())
			continue
		}
		if f.Name() == "expected.go" {
			if err := runTest(directory); err != nil {
				return err
			}
		}
	}
	var errs error
	for _, name := range subDirectories {
		if err := runTests(fmt.Sprintf("%s/%s", directory, name)); err != nil {
			if *exitOnErr {
				return err
			}
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func runTest(directory string) error {
	log.Printf("Running test in %s\n", directory)
	cmd := exec.Command("go", "generate", directory)
	if *verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}
	expected, err := os.ReadFile(fmt.Sprintf("%s/expected.go", directory))
	if err != nil {
		return fmt.Errorf("failed to read expected.go file: %w", err)
	}
	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	var actual []byte
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".gen.go") {
			actual, err = os.ReadFile(fmt.Sprintf("%s/%s", directory, f.Name()))
			if err != nil {
				return fmt.Errorf("failed to read car.gen.go file: %w", err)
			}
			if err := os.Remove(fmt.Sprintf("%s/%s", directory, f.Name())); err != nil {
				return err
			}
			break
		}
	}
	if string(expected) != string(actual) {
		if *verbose {
			log.Printf("expected.go:\n%s\n", string(expected))
			log.Printf("generated file:\n%s\n", string(actual))
		}
		return fmt.Errorf("expected.go and generated file are different")
	}
	cmd = exec.Command("go", "test", "-v", directory)
	if *verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}
	return nil
}
