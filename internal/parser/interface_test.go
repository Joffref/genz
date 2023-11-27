package parser

import (
	"github.com/Joffref/genz/internal/utils"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestParseInterfaceSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode            string
		interfaceName     string
		expectedInterface Interface
	}{
		"basic interface": {
			goCode: `
			package main

			type A interface {}
			`,
			interfaceName: "A",
			expectedInterface: Interface{
				Type:    Type{Name: "main.A", InternalName: "A"},
				Methods: []Method{},
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
			expectedInterface: Interface{
				Type: Type{Name: "main.A", InternalName: "A"},
				Methods: []Method{
					{
						Name:              "Foo",
						Params:            []Type{},
						Returns:           []Type{},
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
			expectedInterface: Interface{
				Type: Type{Name: "main.A", InternalName: "A"},
				Methods: []Method{
					{
						Name:              "Foo",
						Params:            []Type{},
						Returns:           []Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
					},
					{
						Name:              "Bar",
						Params:            []Type{},
						Returns:           []Type{},
						IsPointerReceiver: false,
						IsExported:        true,
						Comments:          []string{},
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

			iface, err := ParseInterface(pkg, tc.interfaceName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(iface, tc.expectedInterface) {
				t.Fatalf("output interface doesn't match expected:\n%s", cmp.Diff(iface, tc.expectedInterface))
			}
		})
	}
}
