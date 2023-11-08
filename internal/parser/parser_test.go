package parser_test

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/Joffref/genz/internal/parser"
	"golang.org/x/tools/go/packages"
)

func TestParseSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode         string
		structName     string
		expectedStruct parser.Struct
	}{
		"basic struct": {
			goCode: `
			package main
			
			type A struct {}
			`,
			structName: "A",
			expectedStruct: parser.Struct{
				Type:       parser.Type("A"),
				Attributes: []parser.Attribute{},
			},
		},
		"struct with one attribute": {
			goCode: `
			package main
			
			type A struct {
				foo string
			}
			`,
			structName: "A",
			expectedStruct: parser.Struct{
				Type: parser.Type("A"),
				Attributes: []parser.Attribute{
					{
						Name:     "foo",
						Type:     "string",
						Comments: []string{},
					},
				},
			},
		},
		"struct with two attributes": {
			goCode: `
			package main
			
			type A struct {
				foo string
				bar uint
			}
			`,
			structName: "A",
			expectedStruct: parser.Struct{
				Type: parser.Type("A"),
				Attributes: []parser.Attribute{
					{
						Name:     "foo",
						Type:     "string",
						Comments: []string{},
					},
					{
						Name:     "bar",
						Type:     "uint",
						Comments: []string{},
					},
				},
			},
		},
		"attribute with doc": {
			goCode: `
			package main
			
			type A struct {
				//comment 1
				//comment 2
				foo string
			}
			`,
			structName: "A",
			expectedStruct: parser.Struct{
				Type: parser.Type("A"),
				Attributes: []parser.Attribute{
					{
						Name:     "foo",
						Type:     "string",
						Comments: []string{"comment 1", "comment 2"},
					},
				},
			},
		},
		"attribute with inline comment": {
			goCode: `
			package main
			
			type A struct {
				foo string // foo
			}
			`,
			structName: "A",
			expectedStruct: parser.Struct{
				Type: parser.Type("A"),
				Attributes: []parser.Attribute{
					{
						Name:     "foo",
						Type:     "string",
						Comments: []string{},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pkg := createPkgWithCode(t, tc.goCode)

			gotStruct, err := parser.Parse(pkg, tc.structName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(gotStruct, tc.expectedStruct) {
				t.Fatalf("expected %s, got %s", tc.expectedStruct, gotStruct)
			}
		})
	}
}

func createPkgWithCode(t *testing.T, goCode string) *packages.Package {
	t.Helper()

	tmp := t.TempDir()
	err := os.WriteFile(path.Join(tmp, "main.go"), []byte(goCode), 0644)
	if err != nil {
		t.Fatalf("failed while writing file: %v", err)
	}

	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax, Tests: false}
	pkgs, err := packages.Load(cfg, path.Join(tmp, "main.go"))
	if err != nil {
		t.Fatalf("failed to load package: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}

	return pkgs[0]
}
