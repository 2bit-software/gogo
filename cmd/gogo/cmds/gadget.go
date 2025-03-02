package cmds

import (
	"fmt"
	"github.com/2bit-software/gogo"
	"github.com/urfave/cli/v2"
)

func GadgetCommand() *cli.Command {
	return &cli.Command{
		Name:            "gadget",
		Usage:           "gogo gadget <function> [args...]",
		Description:     `Run a gogo function. To override GoGo behavior, use environment variables.`,
		SkipFlagParsing: true,
		HideHelpCommand: true, // required so that we can manually capture it in the action
		Action:          gadgetAction,
		Subcommands:     nil,
	}
}

func gadgetAction(ctx *cli.Context) error {
	args := ctx.Args().Slice()
	// build our program arguments
	opts, err := BuildOptions(ctx)
	if err != nil {
		return err
	}

	// first check if the first arg is --help, which means we just want help from the gadget command
	if len(args) >= 1 && args[0] == "--help" {
		err = cli.ShowCommandHelp(ctx, "gadget")
		return err
	}

	// run the command
	err = gogo.Run(opts, args)
	if err != nil && len(args) == 0 {
		fmt.Printf("error: %v\n", err)
		// if we have an error and no arguments, print the help
		err = cli.ShowCommandHelp(ctx, "gadget")
	}
	return err
}
