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
	"github.com/2bit-software/gogo/pkg/sh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var shellCompletionTarget string
var completionInfo string

func init() {
	if err := prepareCommand(rootCmd); err != nil {
		fmt.Printf("error setting flags: %v\n", err)
	}

	// set context of cmd so we can use it to pass data to subcommands
	rootCmd.SetContext(context.Background())

	viper.SetEnvPrefix("GOGO")
	// make sure it sources environment variables as well
	viper.AutomaticEnv()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gogo",
	Short: "A decent JIT-like Go task runner",
	Long:  `Provides a way to generate CLI libraries from a collection of functions, and optionally run them.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// if no args are provided, it returns the complete list of available functions/commands
		// this does NOT expect handling any auto-complete past the first argument, since that *should*
		// be auto-completed by the built binary/function itself. The auto-completion script
		// should detect when the auto-completion is for an argument past the first, request the necessary information
		// from this binary using --autocomplete=<funcName>, and then request auto-completion information from the built
		// binary/function itself.
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		validTargets, _ := gogo.BuildFuncList(gogo.RunOpts{})
		fmt.Println(validTargets)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true, // disable default completion for the root command
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
		version := viper.GetBool("version")
		if version {
			cmd.Printf("%+v\n", Version())
			return nil
		}
		if shellCompletionTarget != "" {
			return gogo.GenerateShellCompletion(shellCompletionTarget, cmd)
		}

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

func prepareCommand(cmd *cobra.Command) error {
	// silence usage on error
	cmd.SilenceUsage = true

	// Whitelist unknown flags, so we can pass them to the subcommands
	cmd.FParseErrWhitelist.UnknownFlags = true

	// set flags
	cmd.Flags().Bool("version", false, "Print the version.")
	cmd.Flags().BoolP("build-global", "b", false, "Build the global cache.")
	cmd.Flags().BoolP("build-local", "l", false, "Build the local cache and exit.")
	cmd.Flags().BoolP("optimize", "o", false, "Optimize the resulting binaries for this run.")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output.")
	cmd.Flags().BoolP("keep-artifacts", "k", false, "Keep the .go files and built binaries.")
	cmd.Flags().BoolP("individual-binaries", "i", false, "Every function gets its own binary, without any subcommands.")
	cmd.Flags().BoolP("disable-cache", "d", false, "Disable cache, forces everything to rebuild.")
	cmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogo.yaml).")
	cmd.Flags().StringVar(&shellCompletionTarget, "shell-completions", "", "Generate shell completions for the given shell, [bash, zsh, fish, powershell].")
	cmd.Flags().StringVar(&completionInfo, "autocomplete", "", "Internal flag used by shell completion to determine where to redirect completion requests.")

	// Bind flags to viper and check for errors
	if err := viper.BindPFlag("VERBOSE", cmd.Flags().Lookup("verbose")); err != nil {
		return err
	}
	if err := viper.BindPFlag("KEEP_ARTIFACTS", cmd.Flags().Lookup("keep-artifacts")); err != nil {
		return err
	}
	if err := viper.BindPFlag("INDIVIDUAL_BINARIES", cmd.Flags().Lookup("individual-binaries")); err != nil {
		return err
	}
	if err := viper.BindPFlag("DISABLE_CACHE", cmd.Flags().Lookup("disable-cache")); err != nil {
		return err
	}
	if err := viper.BindPFlag("OPTIMIZE", cmd.Flags().Lookup("optimize")); err != nil {
		return err
	}
	if err := viper.BindPFlag("BUILD_LOCAL", cmd.Flags().Lookup("build-local")); err != nil {
		return err
	}
	if err := viper.BindPFlag("VERSION", cmd.Flags().Lookup("version")); err != nil {
		return err
	}
	return nil
}

// buildOptions binds the command flags to viper and reads the values from viper
func buildOptions() (gogo.RunOpts, error) {
	// Read the values from viper
	runOpts := gogo.RunOpts{
		Verbose:          viper.GetBool("VERBOSE"),
		GlobalSourceDir:  viper.GetString("GLOBAL_DIR"),
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
