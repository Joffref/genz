package parser

import (
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/pkg/models"
	"golang.org/x/tools/go/packages"
)

// Parse returns a models.ParsedElement from the given *packages.Package and type name.
// It's the main entry point for the different underlying parsers (struct, interface, etc).
func Parse(pkg *packages.Package, typeName string) (models.ParsedElement, error) {
	return parser.Parse(pkg, typeName)
}
