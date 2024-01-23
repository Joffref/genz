# How to write GenZ templates?

GenZ is a template-based generator for Go. A single binary that can be called with the native `go generate` to automate generation of your Go code based on templates.
Thus, a key part of GenZ is writing templates that will be used to generate code.

## Template Syntax

GenZ uses [go templates](https://pkg.go.dev/text/template) to generate code. Thus, you can use all the features of go templates to generate your code.
additionally, we included [sprig](http://masterminds.github.io/sprig) to provide additional functions to the templates (the same as in Helm)

## Template Data

The data injected into the templates can vary depending on the type used to generate the code.

## Template Examples

We provide a few examples of templates in the [examples](../examples) directory.

## Template Testing

We provide a few examples of templates testing in the [examples](../examples) directory or in the [testing.md](../testing.md) file.