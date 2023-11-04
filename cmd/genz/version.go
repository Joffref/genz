package genz

import (
	"flag"
	"github.com/Joffref/genz/internal/command"
	"log"
)

const (
	versionUsage = "print the version"
)

type versionCommand struct {
}

var (
	version = flag.NewFlagSet("version", flag.ExitOnError)
)

func init() {
	version.Usage = func() {
		log.Printf("%s\n", versionUsage)
	}
	command.RegisterCommand("version", versionCommand{})
}

func (v versionCommand) FlagSet() *flag.FlagSet {
	return version
}

func (v versionCommand) Run() error {
	log.Printf("genz version %s\n", Version)
	return nil
}

func (v versionCommand) ValidateArgs() error {
	return nil
}
