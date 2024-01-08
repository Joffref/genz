package generator

import (
	"bytes"
	"fmt"
)

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
