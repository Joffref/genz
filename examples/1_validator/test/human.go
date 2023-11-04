package test

//go:generate genz -type Human -template ../main.tmpl -output human_validator.gen.go
type Human struct {
	Name string
	Age  int
}
