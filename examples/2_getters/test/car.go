package test

//go:generate genz -type Car -template ../getters.tmpl -output car.gen.go
type Car struct {

	//+getter
	model string

	//+getter
	wheels []Wheel

	motorReference string
}

type Wheel struct{}
