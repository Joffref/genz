package main

import (
	"github.com/Joffref/genz/genz/testutils"
	"os"
	"testing"
)

func TestMyCustomGenerator(t *testing.T) {
	typeFile, err := os.ReadFile("./test/car.go")

	parsedElement := testutils.ParseElement(t, string(typeFile), "Car")

	expected, err := os.ReadFile("./test/expected.go")
	if err != nil {
		t.Error(err)
	}

	testutils.IsExpected(t, string(expected), MyCustomGenerator, parsedElement)
}
