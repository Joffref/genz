package command

import "flag"

type Command interface {
	FlagSet() *flag.FlagSet
	Run() error
	ValidateArgs() error
}

var (
	rootCommand Command                // root command
	commands    = map[string]Command{} // available commands name -> Command
)

func RegisterCommand(name string, cmd Command) {
	commands[name] = cmd
}

func SetRootCommand(cmd Command) {
	rootCommand = cmd
}

func Commands() map[string]Command {
	return commands
}

func RootCommand() Command {
	return rootCommand
}
