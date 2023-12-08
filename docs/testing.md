# Test your templates

## Why testing templates?

Testing your templates is a good way to ensure that:
* the generated code is exactly the one you expect
* the generated code does what you expect from it (by running it against unit tests)
It's also a good way catch regressions when you update GenZ or your templates themselves.

## How to test your templates?

GenZ provides a simple way to test your templates using the `genz test` command.

### Test your templates using the CLI

You can test your templates using the CLI with the following command:

```bash
genz test ./path/to/folder/containing/templates
```

## How to write tests for your templates?

GenZ provides a simple way to write tests for your templates. The tests are written in the same language as your templates: Go.
The tests are written in subfolders of your templates folder, as shown in the following example:

```bash
.
├── human.tmpl
└── human_test
    ├── human.go
    ├── expected_test.go # optional
    └── expected.go 
```

With this structure, GenZ will be able to test your template against the expected output and the expected business logic.

### Output testing

GenZ will test the output of your template against the expected output.
The expected output is the content of the `expected.go` file.

### Business logic testing

GenZ will test the business logic of your template against the expected business logic.
The expected business logic is the content of the `expected_test.go` file.

### Example

Let's take the following template as an example:

```mustache
package main

{{ range .Attributes }}
  {{ if has "+getter" .Comments }}{{ $receiverName := substr 0 1 $.Type.InternalName | lower}}
func ({{ $receiverName }} *{{ $.Type.InternalName }}) Get{{ camelcase .Name }}() {{ .Type.InternalName }} {
  return {{ $receiverName }}.{{.Name}}
}
  {{ end }}
{{ end }}
```

We can write the following `car.go` file to test the output of the template (located in `car_test/test_1/car.go`):

```go
package main

//go:generate genz -type Car -template ../../getters.tmpl -output car.gen.go
type Car struct {
	//+getter
	model string
}
```

We can write the following `expected.go` file to test the output of the template (located in `car_test/test_1/expected.go`):

```go
package main

func (c *Car) GetModel() string {
	return c.model
}
```

We can write the following `expected_test.go` file to test the business logic of the template (located in `car_test/test_1/expected_test.go`):

```go
package main
import "testing"
func TestCarGetModel(t *testing.
	c := &Car{}
	if c.GetModel() != c.model {
		t.Errorf("Expected %s, got %s", c.model, c.GetModel()
	}
	c.model = "Ford"
	if c.GetModel() != "Ford" {
		t.Errorf("Expected %s, got %s", c.model, c.GetModel()
	}
	c.model = "Ferrari"
	if c.GetModel() != "Ferrari"
		t.Errorf("Expected %s, got %s", c.model, c.GetModel()
	}
}
```
