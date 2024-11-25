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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/2bit-software/gogo"
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.FParseErrWhitelist.UnknownFlags = true
	runCmd.Flags().BoolP("keep-artifacts", "k", false, "Keep the .go files and built binaries.")
	runCmd.Flags().BoolP("disable-cache", "d", false, "Disable cache, forces everything to rebuild.")
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "func",
	Short: "Run the go function.",
	Long:  `Run the go function.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		funcList, err := gogo.BuildFuncList(gogo.RunOpts{
			Verbose:         false,
			GlobalSourceDir: "",
			GlobalBinDir:    "",
		})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var completions []string
		toCompleteLower := strings.ToLower(toComplete)
		for _, f := range funcList {
			if strings.HasPrefix(strings.ToLower(f.Name), toCompleteLower) {
				completions = append(completions, f.Name)
			}
		}

		// Log the args, toComplete, and completions to /tmp/gogo.log
		//logFile, err := os.OpenFile("/tmp/gogo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		//if err == nil {
		//	defer logFile.Close()
		//	logData := fmt.Sprintf("Args: %v\nToComplete: %s\nCompletions: %v\n", args, toComplete, completions)
		//	logFile.WriteString(logData)
		//}

		return completions, cobra.ShellCompDirectiveDefault
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Get the original args from os.Args
		originalArgs := os.Args[2:] // Skip the program name

		// Parse hidden flags and get remaining args
		gogoArgs, subCmdArgs := gogo.ParseHiddenFlags(originalArgs)

		if len(gogoArgs) == 0 {
			// if there are no remaining args, run the command without any trickery
			return nil
		}

		// in this case we have flags we actually want parsed by gogo, and don't want
		// to manually re-parse them, so we clone the rootCmd without the PreRunE, and
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
		newCmd.Flags().BoolP("keep-artifacts", "k", false, "Keep the .go files and built binaries.")
		newCmd.Flags().BoolP("disable-cache", "d", false, "Disable cache, forces everything to rebuild.")

		if err := viper.BindPFlag("KEEP_ARTIFACTS", newCmd.Flags().Lookup("keep-artifacts")); err != nil {
			return err
		}
		if err := viper.BindPFlag("DISABLE_CACHE", newCmd.Flags().Lookup("disable-cache")); err != nil {
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
