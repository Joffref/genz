package utils

import (
	"fmt"
	"log"
	"os"
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
