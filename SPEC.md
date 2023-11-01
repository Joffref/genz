// Params is the parameters of the template
type Params struct {
    // Type is the name of the type to match
    Type Type // struct
    // Description is the description of the type to match
    Description string
    // Attributes is the list of attributes of the type to match
    Attributes []Attribute
    // Methods is the list of methods of the type to match
    Methods []Method
}

type Type struct {
    // Symbol is the symbol of the type
    Symbol Symbol
    // IsPointer is true if the type is a pointer
    IsPointer bool
    // IsSlice is true if the type is a slice
    IsSlice bool
}

type Symbol struct {
    // Name is the name of the type
    Name string // UUID
    // Package is the package of the type
    Package string // uuid
    // ImportPath is the path of the type
    ImportPath string // github.com/google/uuid
}

// Attribute is the representation of an attribute
type Attribute struct {
    // Name is the name of the attribute
    Name string
    // Description is the description of the attribute
    Description string
    // Type is the type of the attribute
    Type Type
    // Tags is the list of tags of the attribute
    Tags []Tag
    // Keys is the list of keys of the attribute
    Keys map[string]string // //+custom:rulecore
}

// Tag is the representation of a tag
type Tag struct {
    // Name is the name of the tag
    Name string
    // Value is the value of the tag
    Value string
}

// Method is the representation of a method
type Method struct {
    // Name is the name of the method
    Name string
    // Description is the description of the method
    Description string
    // Params is the list of parameters of the method
    Params map[string]Type
    // Returns is the list of returns of the method
    Returns map[string]Type
    // Receivers is the receiver of the method
    Receivers Type
}
