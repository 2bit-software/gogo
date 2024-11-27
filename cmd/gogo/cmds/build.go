// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package cmds

import (
	"fmt"

	"github.com/2bit-software/gogo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the functions.",
	Long: `Pre-build the functions for faster execution. This will build the functions
for either the local or global cache. By default, the local cache is used. The global cache
can be built with the --global flag.

This command bypasses all caches.

You can configure this using the flags, and the .gogoconfig file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// build our program arguments
		opts, err := BuildOptions()
		if err != nil {
			return err
		}

		if opts.BuildGlobalCache {
			err = gogo.BuildGlobal(opts)
			return err
		}

		// run the command
		err = gogo.BuildLocal(opts)
		return err
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolP("global", "g", false, "Build the global function cache.")
	buildCmd.Flags().BoolP("optimize", "o", false, "Optimize the builds.")
	buildCmd.Flags().BoolP("individual-binaries", "i", false, "Each function outputs an individual binary.")
	buildCmd.Flags().BoolP("keep-artifacts", "k", false, "Keep all intermediary artifacts, like the main.gogo.go file")

	if err := viper.BindPFlag("GLOBAL", buildCmd.Flags().Lookup("global")); err != nil {
		fmt.Printf("error setting flags: %v\n", err)
	}
	if err := viper.BindPFlag("KEEP_ARTIFACTS", buildCmd.Flags().Lookup("keep-artifacts")); err != nil {
		fmt.Printf("error setting flags: %v\n", err)
	}
	if err := viper.BindPFlag("INDIVIDUAL_BINARIES", buildCmd.Flags().Lookup("individual-binaries")); err != nil {
		fmt.Printf("error setting flags: %v\n", err)
	}
	if err := viper.BindPFlag("OPTIMIZE", buildCmd.Flags().Lookup("optimize")); err != nil {
		fmt.Printf("error setting flags: %v\n", err)
	}
}
