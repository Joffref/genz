package testutils

import (
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/internal/testutils"
	"github.com/Joffref/genz/pkg/models"
	"testing"
)

// ParseElement parses the given code looking for type and returns the parsed element
// If it fails, it fails the test
func ParseElement(t *testing.T, code string, typeName string) models.ParsedElement {
	p := testutils.CreatePkgWithCode(t, code)
	element, err := parser.Parse(p, typeName)
	if err != nil {
		t.Error(err)
	}
	return element
}
