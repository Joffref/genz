package testutils

import (
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/internal/testutils"
	"github.com/Joffref/genz/pkg/models"
	"testing"
)

// ParseElement parses the given input and returns the parsed element
// If it fails, it fails the test
func ParseElement(t *testing.T, input string, typeName string) models.ParsedElement {
	p := testutils.CreatePkgWithCode(t, input)
	element, err := parser.Parser(p, typeName)
	if err != nil {
		t.Error(err)
	}
	return element
}
