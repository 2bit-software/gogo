package cmds

import (
	"os"

	"github.com/spf13/viper"

	"github.com/2bit-software/gogo"
	"github.com/2bit-software/gogo/pkg/sh"
)

// BuildOptions binds the command flags to viper and reads the values from viper
func BuildOptions() (gogo.RunOpts, error) {
	// Read the values from viper
	runOpts := gogo.RunOpts{
		Verbose:          viper.GetBool("VERBOSE"),
		GlobalSourceDir:  viper.GetString("GLOBAL_SOURCE_DIR"),
		GlobalBinDir:     viper.GetString("GLOBAL_BIN_DIR"),
		BuildLocalCache:  viper.GetBool("BUILD_LOCAL"),
		BuildGlobalCache: viper.GetBool("BUILD_GLOBAL"),
		BuildOpts: gogo.BuildOpts{
			KeepArtifacts:      viper.GetBool("KEEP_ARTIFACTS"),
			IndividualBinaries: viper.GetBool("INDIVIDUAL_BINARIES"),
			DisableCache:       viper.GetBool("DISABLE_CACHE"),
			Optimize:           viper.GetBool("OPTIMIZE"),
			SourceDir:          viper.GetString("BUILD_DIR"),
		},
	}

	width := sh.DetermineWidth(runOpts.Verbose)
	if width > 0 {
		runOpts.ScreenWidth = width
	}

	cwd, err := os.Getwd()
	if err != nil {
		return runOpts, err
	}

	runOpts.OriginalWorkingDir = cwd

	return runOpts, nil
}
