# GenZ

[![Go Reference](https://pkg.go.dev/badge/github.com/Joffref/genz.svg)](https://pkg.go.dev/github.com/Joffref/genz)
[![Go Report Card](https://goreportcard.com/badge/github.com/Joffref/genz)](https://goreportcard.com/report/github.com/Joffref/genz)
[![License](https://img.shields.io/github/license/Joffref/genz)](LICENSE)
![GitHub stars](https://img.shields.io/github/stars/Joffref/genz)
> At this moment, this project is in development. It is not ready for production use.
> Refers to the [Roadmap](#roadmap) section for more information on what has been done.
>
> See the [Contributing](#contributing) section if you would like to help make it better.


<p align="center" width="100%">
    <img width="50%" src="./docs/assets/logo_500_500.png"> 
</p>

GenZ is your all-in-one opinionated code generator for Go. It leverages the power of Go templates to generate code from your own templates.

Thanks to GenZ, you can generate code for your own needs, without having to write a single line of code.


## Installation

You can install GenZ in different ways, depending on your needs.
We strongly recommend using the Go version for production use.

### Using Go - Production

```bash
go install github.com/Joffref/genz
```

### From source - Development

```bash
git clone https://github.com/Joffref/genz.git
cd genz
make install
```

## Usage

Here is a simple example of how to use GenZ.

### Go

Let say you have a `main.tmpl` file in your project, with the following content:
```gotemplate
package main

import (
    "fmt"
)

func (v {{ .Type }}) Validate() error {
{{ range .Attributes }} {{$attribute := .}}
{{ if eq .Type "string" }}{{ range .Comments }}
{{ if eq . "+required" }}
    if v.{{ $attribute.Name }} == "" {
        return fmt.Errorf("attribute '{{ $attribute.Name }}' must be set")
    }
{{ end }}{{ end }}{{ end }}{{ end }}
    return nil
}
```
And the following `human.go` file:
```go
package main

type Human struct {
	// +required
	Firstname string
	// +required
	Lastname string
	Age uint
}
```

You can generate a `human.gen.go` file with the following command:
```bash
genz -type Human -template main.tmpl -output human.gen.go
```

or using the `go generate` command:
```go
//go:generate genz -type Human -template main.tmpl -output human.gen.go
```

The generated file will look like this:
```go
package main

import (
    "fmt"
)

func (v Human) Validate() error {
    if v.Firstname == "" {
        return fmt.Errorf("attribute 'Firstname' must be set")
    }
    if v.Lastname == "" {
        return fmt.Errorf("attribute 'Lastname' must be set")
    }
    return nil
}
```


### CLI
```bash
Usage of genz:
        genz [flags] -type T -template foo.tmpl [directory]
        genz [flags] -type T -template foo.tmpl files... # Must be a single package
Flags:
  -output string
        output file name; default srcdir/<type>.gen.go
  -tags string
        comma-separated list of build tags to apply
  -template string
        go-template file name
  -type string
        comma-separated list of type names; must be set
```

## Contributing

If you would like to contribute to this project, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the Apache 2.0 - see the [LICENSE](LICENSE) file for details.

## Code of Conduct

This project is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Roadmap

- TBD