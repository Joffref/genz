package generator_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Joffref/genz/internal/generator"
	"github.com/Joffref/genz/internal/parser"
	"golang.org/x/tools/go/packages"
)

func TestGenerateErrorParse(t *testing.T) {
	parseFunc := func(pkg *packages.Package, structName string) (interface{}, error) {
		return parser.Struct{}, errors.New("failed to parse package")
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
	parseFunc := func(pkg *packages.Package, structName string) (interface{}, error) {
		return parser.Struct{}, nil
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
	parseFunc := func(pkg *packages.Package, structName string) (interface{}, error) {
		return parser.Struct{}, nil
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
	parseFunc := func(pkg *packages.Package, structName string) (interface{}, error) {
		return parser.Struct{
			Type: parser.Type{
				Name: "TypeName",
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

func TestFormatSuccessInvalidGoCode(t *testing.T) {
	src := generator.Format(*bytes.NewBufferString("package[main\n\n\n"))
	if string(src) != "package[main\n\n\n" {
		t.Fatalf("expected formatted code, got: %q", src)
	}
}

func TestFormatSuccessValidGoCode(t *testing.T) {
	src := generator.Format(*bytes.NewBufferString(" package main\n\n\n"))
	if string(src) != "package main\n" {
		t.Fatalf("expected formatted code, got: %q", src)
	}
}
