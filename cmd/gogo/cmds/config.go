// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package cmds

import (
	"fmt"
	"os"

	"github.com/2bit-software/gogo"
	"github.com/2bit-software/gogo/pkg/sh"
	"github.com/urfave/cli/v2"
)

// BuildOptions creates a RunOpts struct from the CLI context and environment
func BuildOptions(ctx *cli.Context) (gogo.RunOpts, error) {
	runOpts := gogo.RunOpts{
		Verbose:          ctx.Bool("verbose"),
		GlobalSourceDir:  getEnvOrDefault("GOGO_GLOBAL_SOURCE_DIR", ""),
		GlobalBinDir:     getEnvOrDefault("GOGO_GLOBAL_BIN_DIR", ""),
		BuildLocalCache:  ctx.Bool("build-local"),
		BuildGlobalCache: ctx.Bool("global"),
		BuildOpts: gogo.BuildOpts{
			KeepArtifacts:      ctx.Bool("keep-artifacts"),
			IndividualBinaries: ctx.Bool("individual-binaries"),
			DisableCache:       ctx.Bool("disable-cache"),
			Optimize:           ctx.Bool("optimize"),
			SourceDir:          getEnvOrDefault("GOGO_BUILD_DIR", ""),
		},
	}

	width := sh.DetermineWidth(runOpts.Verbose)
	if width > 0 {
		runOpts.ScreenWidth = width
	}

	cwd, err := os.Getwd()
	if err != nil {
		return runOpts, fmt.Errorf("failed to get current working directory: %w", err)
	}

	runOpts.OriginalWorkingDir = cwd

	return runOpts, nil
}

// getEnvOrDefault retrieves an environment variable or returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// AppFlags returns the global flags used across the application
func AppFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Verbose output",
			EnvVars: []string{"GOGO_VERBOSE"},
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
			Name:    "global",
			Aliases: []string{"g"},
			Usage:   "Build global cache",
			EnvVars: []string{"GOGO_BUILD_GLOBAL"},
		},
		&cli.BoolFlag{
			Name:    "keep-artifacts",
			Aliases: []string{"k"},
			Usage:   "Keep the .go files and built binaries",
			EnvVars: []string{"GOGO_KEEP_ARTIFACTS"},
		},
		&cli.BoolFlag{
			Name:    "individual-binaries",
			Aliases: []string{"i"},
			Usage:   "Each function outputs an individual binary",
			EnvVars: []string{"GOGO_INDIVIDUAL_BINARIES"},
		},
		&cli.BoolFlag{
			Name:    "disable-cache",
			Aliases: []string{"d"},
			Usage:   "Disable cache, forces everything to rebuild",
			EnvVars: []string{"GOGO_DISABLE_CACHE"},
		},
		&cli.BoolFlag{
			Name:    "optimize",
			Aliases: []string{"o"},
			Usage:   "Optimize the builds",
			EnvVars: []string{"GOGO_OPTIMIZE"},
		},
		&cli.StringFlag{
			Name:    "build-dir",
			Usage:   "Source directory for builds",
			EnvVars: []string{"GOGO_BUILD_DIR"},
		},
	}
}
