package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/packages"
)

func LoadPackage(patterns []string, tags []string) *packages.Package {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
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

func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func AssertOutputIsEqual(expected, actual []byte, verbose bool) error {
	if string(expected) != string(actual) {
		if verbose {
			log.Printf("expected.go:\n%s\n", string(expected))
			log.Printf("generated file:\n%s\n", string(actual))
		}
		return fmt.Errorf("expected.go and generated file are different")
	}
	return nil
}

// RunCommand runs a command and returns an error if it fails.
// If verbose is true, the command output is printed.
func RunCommand(command []string, verbose bool) error {
	cmd := exec.Command(command[0], command[1:]...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		var commandStr string
		for _, c := range command {
			commandStr += c + " "
		}
		return fmt.Errorf("failed to run %s: %w", commandStr, err)
	}
	return nil
}
