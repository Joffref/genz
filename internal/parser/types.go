package parser

type (
	ParsedType struct {
		Type Type

		Attributes []Attributes
	}

	Attributes struct {
		Name string
		Type Type

		Comments []string
	}

	Type string
)
