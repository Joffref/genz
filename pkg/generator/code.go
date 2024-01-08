package generator

import (
	"bytes"
	"fmt"
	"github.com/Joffref/genz/internal/utils"
	"io"
)

type Code struct {
	writer         io.Writer
	headerComments []string
	packageName    string
	imports        []string
	namedImports   map[string]string
	declarations   []Declaration
}

func NewCode(writer io.Writer, packageName string) *Code {
	return &Code{writer: writer, packageName: packageName}
}

func (c *Code) WithHeaderComments(comments ...string) *Code {
	c.headerComments = append(c.headerComments, comments...)
	return c
}
func (c *Code) WithPackageName(packageName string) *Code {
	c.packageName = packageName
	return c
}

func (c *Code) WithDeclarations(declarations ...Declaration) *Code {
	c.declarations = append(c.declarations, declarations...)
	return c
}

func (c *Code) WithImports(imports ...string) *Code {
	c.imports = append(c.imports, imports...)
	return c
}

func (c *Code) WithNamedImports(namedImports map[string]string) *Code {
	if c.namedImports == nil {
		c.namedImports = make(map[string]string)
	}
	for name, path := range namedImports {
		c.namedImports[name] = path
	}
	return c
}

func (c *Code) Generate() error {
	var buf bytes.Buffer
	if len(c.headerComments) > 0 {
		for _, comment := range c.headerComments {
			buf.WriteString(fmt.Sprintf("// %s\n", comment))
		}
	}
	if c.packageName != "" {
		buf.WriteString(fmt.Sprintf("package %s\n\n", c.packageName))
	} else {
		return fmt.Errorf("package name is required")
	}

	if len(c.imports) > 0 || len(c.namedImports) > 0 {
		buf.WriteString("import (\n")
		for _, importPath := range c.imports {
			buf.WriteString(fmt.Sprintf("\t\"%s\"\n", importPath))
		}
		for name, importPath := range c.namedImports {
			buf.WriteString(fmt.Sprintf("\t%s \"%s\"\n", name, importPath))
		}
		buf.WriteString(")\n\n")
	}

	for _, declaration := range c.declarations {
		buf.WriteString(declaration.Generate())
		buf.WriteString("\n\n")
	}

	_, err := c.writer.Write(utils.Format(buf))
	if err != nil {
		return fmt.Errorf("failed to write generated code: %v", err)
	}
	return nil
}
