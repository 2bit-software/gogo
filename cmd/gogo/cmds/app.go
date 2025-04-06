// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package cmds

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/2bit-software/gogo"
	"github.com/2bit-software/gogo/pkg/gadgets"
)

// rootAction handles the default command behavior
func rootAction(ctx *cli.Context) error {
	if ctx.Bool("version") {
		fmt.Printf("%+v\n", gogo.Version())
		return nil
	}
	if ctx.Bool("build-info") {
		gogo.PrintVersion(os.Stdout)
		return nil
	}

	// build our program arguments
	opts, err := BuildOptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to build options: %w", err)
	}

	count, err := gadgets.ShowFuncList(opts)
	if err != nil {
		return fmt.Errorf("failed to show function list: %w", err)
	}

	if count == 0 {
		err = cli.ShowAppHelp(ctx)
		if err != nil {
			return fmt.Errorf("failed to show app help: %w", err)
		}
	} else {
		fmt.Println("Type 'gogo gadget <function>' to run a function, or `gogo [gadget <function>] --help` for more information.")
	}

	return nil
}

// NewApp creates a new CLI application with all commands and flags configured
func NewApp() *cli.App {
	// Set the version flag and response
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "Print the version",
		EnvVars: []string{"GOGO_VERSION"},
		Action: func(context *cli.Context, b bool) error {
			fmt.Printf("%+v\n", gogo.Version())
			return nil
		},
	}
	app := &cli.App{
		Name:        "gogo",
		Usage:       "A decent JIT-like Go task runner",
		Description: `Provides a way to generate CLI libraries from a collection of functions, and optionally run them.`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Usage:   "Print the version",
				EnvVars: []string{"GOGO_VERSION"},
			},
			&cli.BoolFlag{
				Name:    "build-info",
				Usage:   "Print all build information for GoGo",
				EnvVars: []string{"GOGO_BUILD_INFO"},
			},
			&cli.BoolFlag{
				Name:    "keep-artifacts",
				Aliases: []string{"k"},
				Usage:   "Keep the .go files and built binaries",
				EnvVars: []string{"GOGO_KEEP_ARTIFACTS"},
			},
			&cli.BoolFlag{
				Name:    "disable-cache",
				Aliases: []string{"d"},
				Usage:   "Disable cache, forces everything to rebuild",
				EnvVars: []string{"GOGO_DISABLE_CACHE"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Verbose output",
				EnvVars: []string{"GOGO_VERBOSE"},
			},
			// select the exact folder to use for gogo files
			&cli.StringFlag{
				Name:    "source",
				Usage:   "Choose a source folder for gogo files",
				Aliases: nil,
				EnvVars: []string{"GOGO_SOURCE_DIR"},
			},
		},
		Action: rootAction,
		Commands: []*cli.Command{
			GadgetCommand(),
			BuildCommand(),
			InitCommand(),
		},
	}

	return app
}
