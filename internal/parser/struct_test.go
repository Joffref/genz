package parser

import (
	"github.com/Joffref/genz/internal/utils"
	"go/ast"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseStructSuccess(t *testing.T) {
	testCases := map[string]struct {
		goCode         string
		structName     string
		expectedStruct Element
	}{
		"basic struct": {
			goCode: `
			package main

			type A struct {}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
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
			expectedStruct: Element{
				Type: Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "string", InternalName: "string"},
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
			expectedStruct: Element{
				Type: Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "string", InternalName: "string"},
						Comments: []string{},
					},
					{
						Name:     "bar",
						Type:     Type{Name: "uint", InternalName: "uint"},
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
			expectedStruct: Element{
				Type: Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "string", InternalName: "string"},
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
			expectedStruct: Element{
				Type: Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "string", InternalName: "string"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a slice": {
			goCode: `
			package main

			type B struct {
				foo []string
			}
			`,
			structName: "B",
			expectedStruct: Element{
				Type: Type{Name: "main.B", InternalName: "B"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "[]string", InternalName: "[]string"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo A
			}
			`,
			structName: "B",
			expectedStruct: Element{
				Type: Type{Name: "main.B", InternalName: "B"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "main.A", InternalName: "A"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a slice of named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo []A
			}
			`,
			structName: "B",
			expectedStruct: Element{
				Type: Type{Name: "main.B", InternalName: "B"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "[]main.A", InternalName: "[]A"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a map of named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo map[A]A
			}
			`,
			structName: "B",
			expectedStruct: Element{
				Type: Type{Name: "main.B", InternalName: "B"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "map[main.A]main.A", InternalName: "map[A]A"},
						Comments: []string{},
					},
				},
			},
		},
		"attribute with a struct containing named type": {
			goCode: `
			package main

			type A struct {}
			type B struct {
				foo struct {
					bar []A
					baz string
				}
			}
			`,
			structName: "B",
			expectedStruct: Element{
				Type: Type{Name: "main.B", InternalName: "B"},
				Attributes: []Attribute{
					{
						Name:     "foo",
						Type:     Type{Name: "struct{bar []main.A; baz string}", InternalName: "struct{bar []A; baz string}"},
						Comments: []string{},
					},
				},
			},
		},
		"one empty method, value receiver": {
			goCode: `
			package main

			type A struct {}

			func (a A) foo() {}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
				Methods: []Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []Type{},
						Returns:           []Type{},
						Comments:          []string{},
					},
				},
			},
		},
		"one empty method, pointer receiver": {
			goCode: `
			package main

			type A struct {}

			func (a *A) foo() {}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
				Methods: []Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: true,
						Params:            []Type{},
						Returns:           []Type{},
						Comments:          []string{},
					},
				},
			},
		},
		"one method with 1 param and 1 return, value receiver": {
			goCode: `
			package main

			type A struct {}

			func (a A) foo(a string) int {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
				Methods: []Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []Type{{Name: "string", InternalName: "string"}},
						Returns:           []Type{{Name: "int", InternalName: "int"}},
						Comments:          []string{},
					},
				},
			},
		},
		"one method with 1 param and 1 return, named type": {
			goCode: `
			package main

			type T struct {}
			type A struct {}

			func (a A) foo(a T) T {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
				Methods: []Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []Type{{Name: "main.T", InternalName: "T"}},
						Returns:           []Type{{Name: "main.T", InternalName: "T"}},
						Comments:          []string{},
					},
				},
			},
		},
		"one method with 1 param and 1 return, complex named type": {
			goCode: `
			package main

			type T struct {}
			type A struct {}

			func (a A) foo(a map[T]T) struct{name T} {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
				Methods: []Method{
					{
						Name:              "foo",
						IsExported:        false,
						IsPointerReceiver: false,
						Params:            []Type{{Name: "map[main.T]main.T", InternalName: "map[T]T"}},
						Returns:           []Type{{Name: "struct{name main.T}", InternalName: "struct{name T}"}},
						Comments:          []string{},
					},
				},
			},
		},
		"one exported method with 2 params and 2 returns, pointer receiver": {
			goCode: `
			package main

			type A struct {}

			func (a *A) Foo(a string, b uint) (int, error) {
				return 0
			}
			`,
			structName: "A",
			expectedStruct: Element{
				Type:       Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{},
				Methods: []Method{
					{
						Name:              "Foo",
						IsExported:        true,
						IsPointerReceiver: true,
						Params:            []Type{{Name: "string", InternalName: "string"}, {Name: "uint", InternalName: "uint"}},
						Returns:           []Type{{Name: "int", InternalName: "int"}, {Name: "error", InternalName: "error"}},
						Comments:          []string{},
					},
				},
			},
		},
		"imported type": {
			goCode: `
			package main

			import "github.com/google/uuid"
			
			type A struct {
				foo string
				bar uuid.UUID
				baz map[uuid.UUID]uuid.UUID
			}
			`,
			structName: "A",
			expectedStruct: Element{
				Type: Type{Name: "main.A", InternalName: "A"},
				Attributes: []Attribute{
					{
						Name: "foo",
						Type: Type{
							Name:         "string",
							InternalName: "string",
						},
						Comments: []string{},
					},
					{
						Name: "bar",
						Type: Type{
							Name:         "uuid.UUID",
							InternalName: "UUID",
						},
						Comments: []string{},
					},
					{
						Name: "baz",
						Type: Type{
							Name:         "map[uuid.UUID]uuid.UUID",
							InternalName: "map[UUID]UUID",
						},
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

			pkg := utils.CreatePkgWithCode(t, tc.goCode)

			expr, err := loadAstExpr(pkg, tc.structName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			gotStruct, err := parseStruct(pkg, tc.structName, expr.(*ast.StructType))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(gotStruct, tc.expectedStruct) {
				t.Fatalf("output struct doesn't match expected:\n%s", cmp.Diff(gotStruct, tc.expectedStruct))
			}
		})
	}
}
