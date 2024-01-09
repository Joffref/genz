package parser

import (
	"github.com/Joffref/genz/genz/models"
	"github.com/Joffref/genz/internal/testutils"
	"sort"
	"testing"
)

func TestParsePackageSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode          string
		expectedPackage models.ParsedElement
	}{
		"basic package": {
			goCode: `
			package main
			`,
			expectedPackage: models.ParsedElement{
				PackageName:    "main",
				PackageImports: []string{},
			},
		},
		"package with one import": {
			goCode: `
			package main

			import "fmt"
			`,
			expectedPackage: models.ParsedElement{
				PackageName:    "main",
				PackageImports: []string{"fmt"},
			},
		},
		"package with two imports": {
			goCode: `
			package main

			import (
				"fmt"
				"time"
			)
			`,
			expectedPackage: models.ParsedElement{
				PackageName:    "main",
				PackageImports: []string{"time", "fmt"},
			},
		},
		"package with one import alias": {
			goCode: `
			package main

			import f "fmt"
			`,
			expectedPackage: models.ParsedElement{
				PackageName:    "main",
				PackageImports: []string{"fmt"},
			},
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pkg := testutils.CreatePkgWithCode(t, tc.goCode)

			parsedPackage, err := parsePackage(pkg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if parsedPackage.PackageName != tc.expectedPackage.PackageName {
				t.Fatalf("expected package name %s, got %s", tc.expectedPackage.PackageName, parsedPackage.PackageName)
			}
			if len(parsedPackage.PackageImports) != len(tc.expectedPackage.PackageImports) {
				t.Fatalf("expected %d imports, got %d", len(tc.expectedPackage.PackageImports), len(parsedPackage.PackageImports))
			}
			sort.Strings(parsedPackage.PackageImports)
			sort.Strings(tc.expectedPackage.PackageImports)
			for i := range parsedPackage.PackageImports {
				if parsedPackage.PackageImports[i] != tc.expectedPackage.PackageImports[i] {
					t.Fatalf("expected import %s, got %s", tc.expectedPackage.PackageImports[i], parsedPackage.PackageImports[i])
				}
			}
		})
	}
}
