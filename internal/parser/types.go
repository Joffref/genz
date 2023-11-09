package parser

type (
	Struct struct {
		Type Type

		Attributes []Attribute
	}

	Attribute struct {
		Name string
		Type Type

		Comments []string
	}

	Type string
)
