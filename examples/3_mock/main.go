package main

import (
	"bytes"
	"github.com/Joffref/genz/genz"
	"github.com/Joffref/genz/genz/cli"
	"github.com/Joffref/genz/genz/models"
	"io"
	"text/template"
)

func main() {
	if err := cli.NewCommandFromGenerator("mock-gen", MyCustomGenerator).Execute(); err != nil {
		panic(err)
	}
}

func MyCustomGenerator(buf io.Writer, parsedElement models.ParsedElement) error {
	// TODO: implement

	code := genz.NewCode(buf, parsedElement.PackageName).
		WithImports(parsedElement.PackageImports...)

	mockStruct := genz.Struct("Mock" + parsedElement.Type.InternalName)
	var attributes map[string]string
	for _, method := range parsedElement.Methods {
		attributeName := method.Name + "Func"
		tmpl := "func({{ range $index, $element := .Params }} param{{$index}} {{ .Name }}{{ end }}) {{ range .Returns }}{{ .InternalName }}{{ end }}"
		parse, err := template.New("test").Parse(tmpl)
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		if err := parse.Execute(buf, method); err != nil {
			return err
		}
		attributes[attributeName] = buf.String()
	}
	panic("not implemented")

	return code.WithDeclarations(mockStruct).Generate()
}
