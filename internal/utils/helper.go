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

// RunCommand runs a command and returns an error if it fails.
// If verbose is true, the command output is printed.
func RunCommand(command []string, verbose bool) error {
	cmd := exec.Command(command[0], command[1:]...)
	if verbose {
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
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
