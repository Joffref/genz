package parser

import (
	"github.com/Joffref/genz/internal/utils"
	"golang.org/x/tools/go/packages"
)

// LoadPackage is a helper to load a package
// Note: patterns is a list of patterns to match packages (e.g . ; ./... ; github.com/Joffref/genz/...)
func LoadPackage(patterns []string, buildTags ...string) *packages.Package {
	return utils.LoadPackage(patterns, buildTags)
}
