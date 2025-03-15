// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package cmds

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/2bit-software/gogo/pkg/gadgets"
	"github.com/2bit-software/gogo/pkg/sh"
)

// BuildOptions creates a RunOpts struct from the CLI context and environment
func BuildOptions(ctx *cli.Context) (scripts.RunOpts, error) {
	runOpts := scripts.RunOpts{
		Verbose:          ctx.Bool("verbose"),
		GlobalSourceDir:  getEnvOrDefault("GOGO_GLOBAL_SOURCE_DIR", ""),
		GlobalBinDir:     getEnvOrDefault("GOGO_GLOBAL_BIN_DIR", ""),
		BuildLocalCache:  ctx.Bool("build-local"),
		BuildGlobalCache: ctx.Bool("global"),
		BuildOpts: scripts.BuildOpts{
			KeepArtifacts:  ctx.Bool("keep-artifacts"),
			DisableCache:   ctx.Bool("disable-cache"),
			Optimize:       ctx.Bool("optimize"),
			SourceDir:      getEnvOrDefault("GOGO_BUILD_DIR", ""),
			BinaryFilepath: ctx.String("output"),
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
