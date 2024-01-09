package genz

import (
	"fmt"
	"github.com/Joffref/genz/genz/models"
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/internal/utils"
)

// Parse parses a type from a package
// Note: patterns is a list of patterns to match packages (e.g . ; ./... ; github.com/Joffref/genz/...)
func Parse(typeName string, patterns []string, buildTags ...string) (models.ParsedElement, error) {
	pkg := utils.LoadPackage(patterns, buildTags)
	if pkg == nil {
		return models.ParsedElement{}, fmt.Errorf("could not load package")
	}
	element, err := parser.Parser(pkg, typeName)
	if err != nil {
		return models.ParsedElement{}, err
	}
	return element, err
}
