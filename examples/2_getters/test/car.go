package test

//go:generate go run github.com/Joffref/genz/examples/2_getters -type Car -output car.gen.go
type Car struct {

	//+getter
	Model string

	//+getter
	Wheels []Wheel

	MotorReference string
}

type Wheel struct{}
