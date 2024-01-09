package main

import (
	"github.com/Joffref/genz/genz"
	"github.com/Joffref/genz/genz/cli"
	"github.com/Joffref/genz/genz/models"
	"io"
	"slices"
)

func main() {
	cli.NewCommandFromGenerator("getters", MyCustomGenerator).Execute()
}

func MyCustomGenerator(w io.Writer, parsedElement models.ParsedElement) error {
	code := genz.NewCode(w, parsedElement.PackageName)

	for _, attribute := range parsedElement.Attributes {
		if slices.Contains(attribute.Comments, "+getter") {
			code.WithDeclarations(
				genz.Function("Get"+attribute.Name).
					WithReceiver("v", parsedElement.Type.InternalName, false).
					WithReturns(attribute.Type.InternalName).
					WithBody("return v." + attribute.Name),
			)
		}
	}
	return code.Generate()
}
