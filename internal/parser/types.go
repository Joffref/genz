package parser

type (
	Struct struct {
		Type Type

		Attributes []Attribute
		Methods    []Method
	}

	Attribute struct {
		Name string
		Type Type

		Comments []string
	}

	Method struct {
		Name              string
		Params            []Type
		Returns           []Type
		IsPointerReceiver bool
		IsExported        bool

		Comments []string
	}

	Type struct {
		Name string
	}
)
