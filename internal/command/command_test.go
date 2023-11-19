package command_test

import (
	"flag"
	"reflect"
	"testing"

	"github.com/Joffref/genz/internal/command"
)

type testCmd struct{}

func (t testCmd) FlagSet() *flag.FlagSet { panic("not implemented") }
func (t testCmd) Run() error             { panic("not implemented") }
func (t testCmd) ValidateArgs() error    { panic("not implemented") }

func TestRegisterCommandSuccess(t *testing.T) {
	expected := map[string]command.Command{}
	if !reflect.DeepEqual(command.Commands(), expected) {
		t.Fatalf("Commands doesnt match expected")
	}

	command.RegisterCommand("foo", &testCmd{})
	expected = map[string]command.Command{
		"foo": &testCmd{},
	}
	if !reflect.DeepEqual(command.Commands(), expected) {
		t.Fatalf("Commands doesnt match expected")
	}

	command.RegisterCommand("bar", command.RootCommand())
	expected = map[string]command.Command{
		"foo": &testCmd{},
		"bar": command.RootCommand(),
	}
	if !reflect.DeepEqual(command.Commands(), expected) {
		t.Fatalf("Commands doesnt match expected")
	}

	command.RegisterCommand("bar", &testCmd{})
	expected = map[string]command.Command{
		"foo": &testCmd{},
		"bar": &testCmd{},
	}
	if !reflect.DeepEqual(command.Commands(), expected) {
		t.Fatalf("Commands doesnt match expected")
	}
}

func TestRootCommandSuccess(t *testing.T) {
	if command.RootCommand() != nil {
		t.Fatalf("expected nil, got %s", command.RootCommand())
	}

	command.SetRootCommand(testCmd{})
	if !reflect.DeepEqual(command.RootCommand(), testCmd{}) {
		t.Fatalf("expected testCmd, got %s", command.RootCommand())
	}
}
