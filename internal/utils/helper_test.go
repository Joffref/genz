package utils

import (
	"bytes"
	"testing"
)

func TestFormatSuccessInvalidGoCode(t *testing.T) {
	src := Format(*bytes.NewBufferString("package[main\n\n\n"))
	if string(src) != "package[main\n\n\n" {
		t.Fatalf("expected formatted code, got: %q", src)
	}
}

func TestFormatSuccessValidGoCode(t *testing.T) {
	src := Format(*bytes.NewBufferString(" package main\n\n\n"))
	if string(src) != "package main\n" {
		t.Fatalf("expected formatted code, got: %q", src)
	}
}
