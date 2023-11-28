package parser

import (
	"github.com/Joffref/genz/pkg/models"
	"golang.org/x/tools/go/packages"
)

// parsePackage returns a models.ParsedElement from the given *packages.Package.
// It does not parse the package's elements. Only the package's name and imports are parsed.
func parsePackage(pkg *packages.Package) (models.ParsedElement, error) {
	parsedPackage := models.ParsedElement{
		PackageName:    pkg.Name,
		PackageImports: []string{},
	}

	for i := range pkg.Imports {
		parsedPackage.PackageImports = append(parsedPackage.PackageImports, i)
	}

	return parsedPackage, nil
}
