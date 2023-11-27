package test

//go:generate genz -template ../mock.tmpl -output hello.gen.go -type Hello
type Hello interface {
	SayHelloTo(name string) string
}
