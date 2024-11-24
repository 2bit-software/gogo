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
	"strings"

	"github.com/2bit-software/gogo"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Get the original args from os.Args
		originalArgs := os.Args[1:] // Skip the program name

		// Parse hidden flags and get remaining args
		gogoArgs, subCmdArgs := gogo.ParseHiddenFlags(originalArgs)

		if len(gogoArgs) == 0 {
			// if there are no remaining args, run the command without any trickery
			return nil
		}

		// in this case we have flags we actually want parsed by gogo, and don't want
		// to manually re-parse them, so we clone the rootCmd without the persistenRunE, and
		// run it with the remaining args
		// Create a new root command that will handle the remaining args
		// TODO: we may need to special case this for auto-completion commands!
		newCmd := &cobra.Command{
			Use:               cmd.Use,
			Short:             cmd.Short,
			Long:              cmd.Long,
			RunE:              cmd.RunE,
			PreRunE:           nil, // Important: don't include PreRunE to avoid infinite recursion
			ValidArgsFunction: cmd.ValidArgsFunction,
			CompletionOptions: cmd.CompletionOptions,
		}

		err := prepareCommand(newCmd)
		if err != nil {
			return err
		}

		// Store hidden args in the command's context
		ctx := context.WithValue(cmd.Context(), gogo.HiddenArgsKey{}, subCmdArgs)
		newCmd.SetContext(ctx)

		// Important: Transfer all flag values from original command
		var argsToTransfer []string
		for _, arg := range args {
			if strings.HasPrefix(arg, "--") {
				// it's an actual flag we need to set on the new command
				argsToTransfer = append(argsToTransfer, arg)
			}
		}
		if len(argsToTransfer) > 0 {
			err := newCmd.ParseFlags(argsToTransfer)
			if err != nil {
				return err
			}
		}

		// Set the pointer of the incoming command to the new command
		*cmd = *newCmd

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// extract hidden args from the context
		hiddenArgs := cmd.Context().Value(gogo.HiddenArgsKey{})
		if hiddenArgs != nil {
			args = hiddenArgs.([]string)
			args = hiddenArgs.([]string)
		}

		// build our program arguments
		opts, err := buildOptions()
		if err != nil {
			return err
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

func init() {
	rootCmd.AddCommand(runCmd)
}
