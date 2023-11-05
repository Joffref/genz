package parser

type (
	ParsedType struct {
		Type Type

		Attributes []Attributes
	}

	Attributes struct {
		Name string
		Type Type

		Keys []map[string]string
	}

	Type string
)
