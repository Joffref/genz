package generator_test

import (
	"errors"
	"github.com/Joffref/genz/pkg/models"
	"strings"
	"testing"

	"github.com/Joffref/genz/internal/generator"
	"golang.org/x/tools/go/packages"
)

func TestGenerateErrorParse(t *testing.T) {
	parseFunc := func(pkg *packages.Package, structName string) (models.ParsedElement, error) {
		return models.ParsedElement{}, errors.New("failed to parse package")
	}
	_, err := generator.Generate(nil, "template", "typeName", parseFunc)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to inspect package: failed to parse package" {
		t.Fatalf("error does not match expected: %v", err)
	}
}

func TestGenerateErrorTemplateParse(t *testing.T) {
	parseFunc := func(pkg *packages.Package, structName string) (models.ParsedElement, error) {
		return models.ParsedElement{}, nil
	}
	_, err := generator.Generate(nil, "{{", "typeName", parseFunc)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse template") {
		t.Fatalf("expected string not found in error: %v", err)
	}
}

func TestGenerateErrorTemplateExecute(t *testing.T) {
	parseFunc := func(pkg *packages.Package, structName string) (models.ParsedElement, error) {
		return models.ParsedElement{}, nil
	}
	_, err := generator.Generate(nil, "{{.Foo}}", "typeName", parseFunc)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to execute template") {
		t.Fatalf("expected string not found in error: %v", err)
	}
}

func TestGenerateSuccess(t *testing.T) {
	parseFunc := func(pkg *packages.Package, structName string) (models.ParsedElement, error) {
		return models.ParsedElement{
			Element: models.Element{
				Type: models.Type{
					Name: "TypeName",
				},
			},
		}, nil
	}
	buf, err := generator.Generate(nil, "{{.Type.Name}}", "TypeName", parseFunc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "TypeName" {
		t.Fatalf("expected TypeName, got %s", buf.String())
	}
}
