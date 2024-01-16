package models

type (
	// ParsedElement represents a struct or an interface with its package name and its imports.
	// It is the struct used by Genz to fill the templates. It is the result of the parsing of your type.
	ParsedElement struct {
		// PackageName is the name of the package of the parsed element.
		// e.g. "package foo" => "foo"
		PackageName string
		// List of the imports of the package of the parsed element.
		// e.g. ["github.com/google/uuid", "time"]
		PackageImports []string

		// See Element for more details.
		// Note: this inlined, so you can access directly to the fields as if it was an
		// ParsedElement attribute. e.g. {{ .Type.Name }}
		Element
	}

	// Element represents a struct or an interface.
	Element struct {
		// See Type for more details.
		Type Type

		// List of the attributes of the struct. Empty if the parsed element is an interface.
		// See Attribute for more details.
		Attributes []Attribute

		// List of the methods of the struct or the interface.
		// See Method for more details.
		Methods []Method
	}

	// Attribute represents a struct's attribute.
	// It is not used for interfaces.
	Attribute struct {
		// Name of the attribute. e.g. "Foo" for "Foo string"
		Name string
		// See Type for more details.
		Type Type

		// List of the comments of the attribute.
		// Only upper comments are parsed. No inline comments.
		Comments []string

		// Tags of the attribute. e.g. `json:"foo,omitempty"`
		// The map key is the tag name, the map value is the tag value.
		// e.g. `json:"foo,omitempty"` => map[string]string{"json": "foo,omitempty"}
		Tags map[string]string
	}

	// Method represents a struct's method or an interface's method.
	Method struct {
		// Name of the method. e.g. "Foo" for "Foo() string"s
		Name string
		// Maps of the parameters of the method. Empty if the method has no parameter.
		// The map key is the parameter name, the map value is the parameter type.
		// See Type for more details.
		Params map[string]Type
		// List of the return values of the method. Empty if the method has no return value.
		// See Type for more details.
		Returns []Type
		// IsPointerReceiver is true if the method is a pointer receiver.
		// Always false for interfaces.
		IsPointerReceiver bool
		// IsExported is true if the method is exported.
		IsExported bool

		// List of the comments of the method.
		// Only upper comments are parsed. No inline or in the method's body comments.
		Comments []string
	}

	// Type represents a generic type among structs, interfaces, attributes, methods, etc.
	// It is used to represent various types such as "string", "time.Time", "uuid.UUID", etc.
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
)
