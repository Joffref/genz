# How to write GenZ templates?

GenZ is a template-based generator for Go. A single binary that can be called with the native `go generate` to automate generation of your Go code based on templates.
Thus, a key part of GenZ is writing templates that will be used to generate code.

## Template Syntax

GenZ uses [go templates](https://pkg.go.dev/text/template) to generate code. Thus, you can use all the features of go templates to generate your code.
Plus, we use [sprig](http://masterminds.github.io/sprig) to provide additional functions to the templates.

## Template Data

The data injected into the templates can vary depending on the type used to generate the code. Thus, we provide a cheatsheet for each type.
You can find the cheatsheets in the [input_cheatsheet.md](./input_cheatsheet.md) file.

## Template Examples

We provide a few examples of templates in the [examples](../examples) directory.

## Template Testing

We provide a few examples of templates testing in the [examples](../examples) directory or in the [testing.md](../testing.md) file.