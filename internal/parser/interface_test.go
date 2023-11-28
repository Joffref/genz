package parser

import (
	"github.com/Joffref/genz/internal/utils"
	"github.com/Joffref/genz/pkg/models"
	"github.com/google/go-cmp/cmp"
	"go/ast"
	"reflect"
	"testing"
)

func TestParseInterfaceSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode            string
		interfaceName     string
		expectedInterface models.Element
	}{
		"basic interface": {
			goCode: `
			package main

			type A interface {}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type:    models.Type{Name: "main.A", InternalName: "A"},
				Methods: nil,
			},
		},
		"interface with one method": {
			goCode: `
			package main

			type A interface {
				Foo()	
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with two methods": {
			goCode: `
			package main

			type A interface {
				Foo()
				Bar()
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Bar",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with one method with comments": {
			goCode: `
			package main

			type A interface {
				//Foo does something
				Foo()
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{"Foo does something"},
					},
				},
			},
		},
		"interface with one method with params": {
			goCode: `
			package main
			
			type A interface {
				Foo(a int, b string)
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with one method with returns": {
			goCode: `
			package main

			type A interface {
				Foo() (int, string)
			}
			`,
			interfaceName: "A",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.A", InternalName: "A"},
				Methods: []models.Method{
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
				},
			},
		},
		"interface with a sub interface": {
			goCode: `	
			package main
		
			type A interface {
				Foo() (int, string)
			}
			
			type B interface {
				// A is a sub interface
				A
				Bar()
			}
			`,
			interfaceName: "B",
			expectedInterface: models.Element{
				Type: models.Type{Name: "main.B", InternalName: "B"},
				Methods: []models.Method{
					{
						Name:              "Bar",
						Params:            []models.Type{},
						Returns:           []models.Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
					{
						Name:              "Foo",
						Params:            []models.Type{},
						Returns:           []models.Type{{Name: "int", InternalName: "int"}, {Name: "string", InternalName: "string"}},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{" A is a sub interface"},
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pkg := utils.CreatePkgWithCode(t, tc.goCode)

			expr, err := loadAstExpr(pkg, tc.interfaceName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			iface, err := parseInterface(pkg, tc.interfaceName, expr.(*ast.InterfaceType))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(iface, tc.expectedInterface) {
				t.Fatalf("output interface doesn't match expected:\n%s", cmp.Diff(iface, tc.expectedInterface))
			}
		})
	}
}
