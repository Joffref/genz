package parser

import "golang.org/x/tools/go/packages"

func parsePackage(pkg *packages.Package) (ParsedElement, error) {
	parsedPackage := ParsedElement{
		PackageName:    pkg.Name,
		PackageImports: []string{},
	}

	for i := range pkg.Imports {
		parsedPackage.PackageImports = append(parsedPackage.PackageImports, i)
	}

	return parsedPackage, nil
}
