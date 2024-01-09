package test

//go:generate go run github.com/Joffref/genz/examples/1_validator -type Human -output human.gen.go
type Human struct {
	//+startsWithCapital
	Firstname string

	//+required
	//+startsWithCapital
	Lastname string

	//+>18,<99
	Age uint
}
