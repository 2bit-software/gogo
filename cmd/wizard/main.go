// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package main

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/scripts"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/2bit-software/gogo"
	"github.com/2bit-software/gogo/cmd/gogo/cmds"
)

func main() {
	app := createApp()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// createApp constructs and configures the CLI application
func createApp() *cli.App {
	return &cli.App{
		Name:    "wizard",
		Usage:   "A mage-compatible GoGo wizard",
		Version: gogo.Version().Version,
		Description: `An entrypoint you can symlink from mage to which will use GoGo instead. 
It does not support most of the shell enhancements that GoGo provides, but it does 
support all of the function capabilities.`,
		Action:          executeGoGoCommand,
		HideHelpCommand: true,
		SkipFlagParsing: true,
	}
}

// executeGoGoCommand handles the main program execution flow
func executeGoGoCommand(ctx *cli.Context) error {
	// Build our program options
	opts, err := cmds.BuildOptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to build options: %w", err)
	}

	args := ctx.Args().Slice()

	if len(args) == 0 {
		return handleFunctionListDisplay(ctx, opts)
	}

	// selectively show help
	if len(args) == 1 && args[0] == "--help" {
		return cli.ShowAppHelp(ctx)
	}
	if len(args) == 1 && (args[0] == "-v" || args[0] == "--version") {
		cli.ShowVersion(ctx)
		return nil
	}

	// Execute the GoGo command
	err = scripts.Run(opts, args)

	// Return any error from execution
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// handleFunctionListDisplay shows the available functions or help
func handleFunctionListDisplay(ctx *cli.Context, opts scripts.RunOpts) error {
	count, err := scripts.ShowFuncList(opts)
	if err != nil {
		return fmt.Errorf("failed to show function list: %w", err)
	}

	if count == 0 {
		fmt.Printf("No functions found.\n")
		return cli.ShowAppHelp(ctx)
	} else {
		fmt.Println("Type 'wizard <function>' to run a function, or `wizard --help` for more information.")
	}
	return nil
}
