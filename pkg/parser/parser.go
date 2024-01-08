package parser

import (
	"fmt"
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/internal/utils"
	"github.com/Joffref/genz/pkg/models"
)

// Parse parses a type from a package
// Note: patterns is a list of patterns to match packages
func Parse(typeName string, patterns []string, buildTags ...string) (models.ParsedElement, error) {
	pkg := utils.LoadPackage(patterns, buildTags)
	if pkg == nil {
		return models.ParsedElement{}, fmt.Errorf("could not load package")
	}
	fmt.Println("pkg:", pkg.PkgPath)
	element, err := parser.Parser(pkg, typeName)
	if err != nil {
		return models.ParsedElement{}, err
	}
	return element, err
}
