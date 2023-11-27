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

		// Name of the type from outside its package
		// Example `uuid.UUID` or `time.Time`
		// Use this variable if you generate code outside the package of that type
		// Do not forget to add corresponding imports in your template
		Name string

		// Name of the type from inside its package
		// Example `UUID` or `Time`
		// Use this variable if you generate code inside the package of that type
		InternalName string
	}

	Interface struct {
		Type Type

		Methods  []Method
		Comments []string
	}
)
