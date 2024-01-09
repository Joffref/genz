package genz

import (
	"bytes"
	"fmt"
)

type Declaration interface {
	Generate() string
}

type FuncDecl struct {
	name          string
	inputs        map[string]string
	returns       []string
	receiverName  string
	receiverType  string
	isPtrReceiver bool
	isSignature   bool
	body          string
}

func Function(name string) *FuncDecl {
	return &FuncDecl{
		name: name,
	}
}

func (f *FuncDecl) WithReceiver(receiverName string, receiverType string, isPtrReceiver bool) *FuncDecl {
	f.isPtrReceiver = isPtrReceiver
	f.receiverType = receiverType
	f.receiverName = receiverName
	return f
}

func (f *FuncDecl) WithInputs(inputs map[string]string) *FuncDecl {
	if f.inputs == nil {
		f.inputs = make(map[string]string)
	}
	for inputName, inputType := range inputs {
		f.inputs[inputName] = inputType
	}
	return f
}

func (f *FuncDecl) WithReturns(returns ...string) *FuncDecl {
	f.returns = returns
	return f
}

func (f *FuncDecl) WithBody(body string) *FuncDecl {
	f.body = body
	return f
}

func (f *FuncDecl) SignatureOnly() *FuncDecl {
	f.body = ""
	f.isSignature = true
	return f
}

func (f *FuncDecl) Generate() string {
	var buf bytes.Buffer
	buf.WriteString("func ")
	if f.receiverName != "" {
		buf.WriteString("(")
		buf.WriteString(f.receiverName)
		if f.isPtrReceiver {
			buf.WriteString("*")
		}
		buf.WriteString(" ")
		buf.WriteString(f.receiverType)
		buf.WriteString(")")
	}
	buf.WriteString(" ")
	buf.WriteString(f.name)
	buf.WriteString("(")
	for inputName, inputType := range f.inputs {
		buf.WriteString(inputName)
		buf.WriteString(" ")
		buf.WriteString(inputType)
		buf.WriteString(", ")
	}
	buf.WriteString(") ")
	if len(f.returns) > 0 {
		buf.WriteString("(")
		for _, returnType := range f.returns {
			buf.WriteString(returnType)
			buf.WriteString(", ")
		}
		buf.WriteString(")")
	}
	if f.isSignature {
		buf.WriteString("\n")
	} else {
		buf.WriteString(" {")
		if f.body != "" {
			buf.WriteString("\n")
			buf.WriteString(f.body)
			buf.WriteString("\n")
		}
		buf.WriteString("}")
	}
	return buf.String()
}

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

type StructDecl struct {
	name   string
	fields map[string]struct {
		fieldName     string
		fieldType     string
		inlineComment string
		comments      []string
		tags          map[string]string
	}
}

func Struct(name string) *StructDecl {
	return &StructDecl{
		name: name,
		fields: make(map[string]struct {
			fieldName     string
			fieldType     string
			inlineComment string
			comments      []string
			tags          map[string]string
		}),
	}
}

func (s *StructDecl) WithField(fieldName string, fieldType string) *StructDecl {
	s.fields[fieldName] = struct {
		fieldName     string
		fieldType     string
		inlineComment string
		comments      []string
		tags          map[string]string
	}{fieldName: fieldName, fieldType: fieldType}
	return s
}

func (s *StructDecl) WithFieldInlineComment(fieldName string, fieldType string, inlineComment string) *StructDecl {
	s.fields[fieldName] = struct {
		fieldName     string
		fieldType     string
		inlineComment string
		comments      []string
		tags          map[string]string
	}{fieldName: fieldName, fieldType: fieldType, inlineComment: inlineComment}
	return s
}

func (s *StructDecl) WithFieldComments(fieldName string, fieldType string, comments []string) *StructDecl {
	s.fields[fieldName] = struct {
		fieldName     string
		fieldType     string
		inlineComment string
		comments      []string
		tags          map[string]string
	}{fieldName: fieldName, fieldType: fieldType, comments: comments}
	return s
}

func (s *StructDecl) WithFieldTags(fieldName string, fieldType string, tags map[string]string) *StructDecl {
	s.fields[fieldName] = struct {
		fieldName     string
		fieldType     string
		inlineComment string
		comments      []string
		tags          map[string]string
	}{fieldName: fieldName, fieldType: fieldType, tags: tags}
	return s
}

func (s *StructDecl) Generate() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("type %s struct {\n", s.name))
	for _, field := range s.fields {
		if field.inlineComment != "" {
			buf.WriteString(fmt.Sprintf("\t%s %s // %s\n", field.fieldName, field.fieldType, field.inlineComment))
		} else if len(field.comments) > 0 {
			buf.WriteString(fmt.Sprintf("\t%s %s // %s\n", field.fieldName, field.fieldType, field.comments[0]))
			for _, comment := range field.comments[1:] {
				buf.WriteString(fmt.Sprintf("\t// %s\n", comment))
			}
		} else {
			buf.WriteString(fmt.Sprintf("\t%s %s\n", field.fieldName, field.fieldType))
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}
