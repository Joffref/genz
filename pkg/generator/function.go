package generator

import (
	"bytes"
)

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
