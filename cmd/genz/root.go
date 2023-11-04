//nolint:unused
package genz

import (
	"github.com/Joffref/genz/internal/command"
	"os"
)

const (
	// Version is the current version of genz
	Version = "0.0.1"
)

func Execute() error {
	if len(os.Args) == 1 {
		command.RootCommand().FlagSet().Usage()
		return nil
	}
	for name, cmd := range command.Commands() {
		if name == os.Args[1] {
			if err := cmd.FlagSet().Parse(os.Args[2:]); err != nil {
				return err
			}
			if err := cmd.ValidateArgs(); err != nil {
				return err
			}
			return cmd.Run()
		}
	}
	if err := command.RootCommand().FlagSet().Parse(os.Args[1:]); err != nil {
		return err
	}
	if err := command.RootCommand().ValidateArgs(); err != nil {
		return err
	}
	return command.RootCommand().Run()
}
