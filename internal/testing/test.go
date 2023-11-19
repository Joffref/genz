package testing

import (
	"errors"
	"fmt"
	"github.com/Joffref/genz/internal/utils"
	"os"
	"path"
	"strings"
)

func RunTests(directory string, verbose, exitOnError bool) error {
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
			if err := runTest(directory, verbose); err != nil {
				return err
			}
		}
	}
	var errs error
	for _, name := range subDirectories {
		if err := RunTests(path.Join(directory, name), verbose, exitOnError); err != nil {
			if exitOnError {
				return err
			}
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func runTest(directory string, verbose bool) error {
	fmt.Println("running test in", directory)
	if err := utils.RunCommand([]string{"go", "generate", fmt.Sprintf("./%s", directory)}, verbose); err != nil {
		return err
	}
	expected, err := os.ReadFile(path.Join(directory, "expected.go"))
	if err != nil {
		return fmt.Errorf("failed to read expected.go file: %w", err)
	}
	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	if len(dir) == 0 {
		return fmt.Errorf("no file in directory %s", directory)
	}
	var generatedFiles []string
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".gen.go") {
			generatedFiles = append(generatedFiles, f.Name())
		}
	}
	if len(generatedFiles) == 0 {
		return fmt.Errorf("no generated file in directory %s\nSpecify output file with -output flag in genz command", directory)
	}
	if len(generatedFiles) > 1 {
		var errs error
		if err := cleanUpGeneratedFiles(directory, generatedFiles); err != nil {
			errs = errors.Join(errs, err)
		}
		errs = errors.Join(errs, fmt.Errorf("too many generated files in directory %s\nAt the moment, GenZ only supports one generated file per directory for test ", directory))
		return errs
	}
	actual, err := os.ReadFile(path.Join(directory, generatedFiles[0]))
	if err != nil {
		return fmt.Errorf("failed to read generated file: %w", err)
	}
	if err := cleanUpGeneratedFiles(directory, generatedFiles); err != nil {
		return err
	}
	if err := assertOutputIsEqual(expected, actual, verbose); err != nil {
		return err
	}
	return utils.RunCommand([]string{"go", "test", "-v", fmt.Sprintf("./%s", directory)}, verbose)
}
