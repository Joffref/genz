package test

//go:generate genz -type Human -generator ../main.so:MyCustomGenerator -output human_validator.gen.go
type Human struct {
	//+startsWithCapital
	Firstname string

	//+required
	//+startsWithCapital
	Lastname string

	//+>18,<99
	Age uint
}
