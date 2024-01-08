package generator

import (
	"bytes"
	"fmt"
)

type InterfaceDecl struct {
	name     string
	comments []string
	methods  map[string]struct {
		methodName string
		methodDecl *FuncDecl
	}
}

func Interface(name string) *InterfaceDecl {
	return &InterfaceDecl{
		name: name,
		methods: make(map[string]struct {
			methodName string
			methodDecl *FuncDecl
		}),
	}
}

func (i *InterfaceDecl) WithComments(comments ...string) *InterfaceDecl {
	i.comments = append(i.comments, comments...)
	return i
}

func (i *InterfaceDecl) WithMethod(methodName string, methodDecl *FuncDecl) *InterfaceDecl {
	methodDecl.SignatureOnly()
	i.methods[methodName] = struct {
		methodName string
		methodDecl *FuncDecl
	}{methodName: methodName, methodDecl: methodDecl}
	return i
}

func (i *InterfaceDecl) Generate() string {
	var buf bytes.Buffer
	if len(i.comments) > 0 {
		for _, comment := range i.comments {
			buf.WriteString(fmt.Sprintf("// %s\n", comment))
		}
	}
	buf.WriteString(fmt.Sprintf("type %s interface {\n", i.name))
	for _, method := range i.methods {
		buf.WriteString(fmt.Sprintf("%s\n", method.methodDecl.Generate()))
	}
	buf.WriteString("}\n")
	return buf.String()
}
