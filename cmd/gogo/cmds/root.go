// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/2bit-software/gogo"
)

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// set context of cmd so we can use it to pass data to subcommands
	rootCmd.SetContext(context.Background())

	// set flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output.")
	rootCmd.Flags().Bool("version", false, "Print the version.")
	rootCmd.Flags().BoolP("keep-artifacts", "k", false, "Keep the .go files and built binaries.")
	rootCmd.Flags().BoolP("disable-cache", "d", false, "Disable cache, forces everything to rebuild.")
	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogo.yaml).")

	if err := viper.BindPFlag("VERBOSE", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Printf("error setting flags: %v\n", err)
	}

	viper.SetEnvPrefix("GOGO")
	// make sure it sources environment variables as well
	viper.AutomaticEnv()

	// Whitelist unknown flags, so we can pass them to the subcommands
	rootCmd.FParseErrWhitelist.UnknownFlags = true
}

var rootCmd = &cobra.Command{
	Use:   "gogo",
	Short: "A decent JIT-like Go task runner",
	Long:  `Provides a way to generate CLI libraries from a collection of functions, and optionally run them.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Bind flags to viper and check for errors
		if err := viper.BindPFlag("KEEP_ARTIFACTS", cmd.Flags().Lookup("keep-artifacts")); err != nil {
			return err
		}
		if err := viper.BindPFlag("DISABLE_CACHE", cmd.Flags().Lookup("disable-cache")); err != nil {
			return err
		}

		version := cmd.Flags().Changed("version")
		if version {
			cmd.Printf("%+v\n", Version())
			return nil
		}
		// try listing the functions
		if len(args) == 0 {
			// build our program arguments
			opts, err := BuildOptions()
			if err != nil {
				return err
			}

			count, err := gogo.ShowFuncList(opts)
			if err != nil {
				return err
			}
			if count == 0 {
				_ = cmd.Help()
			} else {
				fmt.Println("Type 'gogo func <function>' to run a function, or `gogo --help` for more information.")
			}
		}
		return nil
	},
}
