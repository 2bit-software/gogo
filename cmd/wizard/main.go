// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/2bit-software/gogo"
	"github.com/2bit-software/gogo/cmd/gogo/cmds"
)

func init() {
	// we don't support command line arguments in this version,
	// but still support environment variables
	viper.SetEnvPrefix("GOGO")
	viper.AutomaticEnv()
}

func main() {
	err := runCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "gogo_wizard",
	Short: "A mage-compatible GoGo wizard",
	Long: `An entrypoint you can symlink from mage to which will use GoGo instead. It does not support most of the shell enhancements that GoGo provides, but it does support all of the function
capabilities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// build our program arguments
		opts, err := cmds.BuildOptions()
		if err != nil {
			return err
		}

		// if this is a version request
		if viper.GetBool("VERSION") {
			cmd.Printf("%+v\n", cmds.Version())
			return nil
		}

		// run the command
		err = gogo.Run(opts, args)
		if err != nil && len(args) == 0 {
			// if we have an error and no arguments, print the help
			_ = cmd.Help()
		}
		// if there's an error, add an extra newline before starting to print
		if err != nil {
			fmt.Println("-")
		}
		return err
	},
}
