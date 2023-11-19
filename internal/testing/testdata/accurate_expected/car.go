package main

//go:generate genz -type Car -template ../getters.tmpl -output car.gen.go
type Car struct {
	model string
}
