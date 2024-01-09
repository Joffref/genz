package testutils

import (
	"bytes"
	"github.com/Joffref/genz/genz"
	"github.com/Joffref/genz/genz/models"
	"github.com/Joffref/genz/internal/parser"
	"github.com/Joffref/genz/internal/testutils"
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

// IsExpected checks that the generator generates the expected code
// If not, it fails the test
func IsExpected(t *testing.T, expect string, generator genz.Generator, parsedElement models.ParsedElement) {
	var buffer bytes.Buffer
	err := generator(&buffer, parsedElement)
	if err != nil {
		t.Error(err)
	}
	if buffer.String() != expect {
		t.Errorf("expected %s, got %s", expect, buffer.String())
	}
}
