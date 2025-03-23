// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/2bit-software/gogo/pkg/gadgets"
)

// buildAction handles the main logic for the build command
func buildAction(ctx *cli.Context) error {
	opts, err := BuildOptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to build options: %w", err)
	}

	switch {
	case ctx.Bool("global"):
		return buildGlobalCache(opts)
	case ctx.Bool("gen-only"):
		return generateFilesOnly(opts)
	default:
		return buildLocalCache(opts)
	}
}

// buildGlobalCache handles building the global cache
func buildGlobalCache(opts gadgets.RunOpts) error {
	if err := gadgets.BuildGlobal(opts); err != nil {
		return fmt.Errorf("failed to build global cache: %w", err)
	}
	return nil
}

// generateFilesOnly handles generating Go files without building
func generateFilesOnly(opts gadgets.RunOpts) error {
	fmt.Println("Generating go files only.")
	if err := gadgets.GenerateMainFile(opts); err != nil {
		return fmt.Errorf("failed to generate main file: %w", err)
	}
	return nil
}

// buildLocalCache handles building the local cache
func buildLocalCache(opts gadgets.RunOpts) error {
	if err := gadgets.BuildLocal(opts); err != nil {
		return fmt.Errorf("failed to build local cache: %w", err)
	}
	return nil
}

// BuildCommand creates the build command which pre-builds functions for faster execution.
// It supports both local and global caching with various optimization options.
func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build the functions",
		Description: `Pre-build the functions for faster execution. This will build the functions
for either the local or global cache. By default, the local cache is used. The global cache
can be built with the --global flag.

This command bypasses all caches.

You can configure this using the flags, and the .gogoconfig file.`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "gen-only",
				Usage:   "Generate the go files only",
				EnvVars: []string{"GOGO_GEN_ONLY"},
			},
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "Build the global function cache",
				EnvVars: []string{"GOGO_GLOBAL"},
			},
			&cli.BoolFlag{
				Name:    "optimize",
				Usage:   "Optimize the builds",
				EnvVars: []string{"GOGO_OPTIMIZE"},
			},
			&cli.BoolFlag{
				Name:    "keep-artifacts",
				Aliases: []string{"k"},
				Usage:   "Keep all intermediary artifacts, like the main.gogo.go file",
				EnvVars: []string{"GOGO_KEEP_ARTIFACTS"},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output binary to given file",
				EnvVars: []string{"GOGO_OUTPUT"},
			},
			&cli.StringFlag{
				Name:    "global-source-dir",
				Usage:   "Global source directory",
				EnvVars: []string{"GOGO_GLOBAL_SOURCE_DIR"},
			},
			&cli.StringFlag{
				Name:    "global-bin-dir",
				Usage:   "Global binary directory",
				EnvVars: []string{"GOGO_GLOBAL_BIN_DIR"},
			},
			&cli.BoolFlag{
				Name:    "build-local",
				Usage:   "Build local cache",
				EnvVars: []string{"GOGO_BUILD_LOCAL"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				EnvVars: []string{"GOGO_VERBOSE"},
			},
		},
		Action: buildAction,
	}
}
