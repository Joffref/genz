package validator

//go:generate genz -type Human -template ../main.tmpl -output human_validator.gen.go
type Human struct {

	//+validator=optional
	Name string

	//+validator=>18,<99
	Age uint
}
