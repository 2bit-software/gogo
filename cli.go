package gogo

import (
	"github.com/urfave/cli/v2"
)

// These re-export some of the types from urfave/cli/v2
// so that the sub commands only need to import gogo
// and we cna potentially swap out the underlying cli and/or
// hook into it to make the subcommands smarter

type App = cli.App

// Commands
type Command = cli.Command
type CliContext = cli.Context

// Flags
type Flag = cli.Flag
type BoolFlag = cli.BoolFlag
type StringFlag = cli.StringFlag
type IntFlag = cli.IntFlag

// VersionFlag prints the version for the application
var VersionFlag Flag = &BoolFlag{
	Name:               "version",
	Aliases:            []string{"v"},
	Usage:              "print the version",
	DisableDefaultText: true,
}
