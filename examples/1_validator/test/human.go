package test

//go:generate genz -type Human -template ../main.tmpl -output human_validator.gen.go
type Human struct {
	//+startsWithCapital
	Firstname string

	//+optional
	//+startsWithCapital
	Lastname string

	//+>18,<99
	Age uint
}
