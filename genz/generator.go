package genz

import (
	"github.com/Joffref/genz/genz/models"
	"io"
)

// Generator is a function that generates code from a parsed element.
// Generator is user-defined and can be used to generate any Go code.
type Generator func(buffer io.Writer, parsedElement models.ParsedElement) error
