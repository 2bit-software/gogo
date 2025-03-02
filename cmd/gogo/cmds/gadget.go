package cmds

import (
	"github.com/2bit-software/gogo"
	"github.com/urfave/cli/v2"
)

func GadgetCommand() *cli.Command {
	return &cli.Command{
		Name:            "gadget",
		Usage:           "gogo gadget <function> [args...]",
		Description:     `Run a gogo function. To override GoGo behavior, use environment variables.`,
		SkipFlagParsing: true,
		Action:          gadgetAction,
	}
}

func gadgetAction(ctx *cli.Context) error {
	args := ctx.Args().Tail()
	// build our program arguments
	opts, err := BuildOptions(ctx)
	if err != nil {
		return err
	}

	// run the command
	err = gogo.Run(opts, args)
	if err != nil && len(args) == 0 {
		// if we have an error and no arguments, print the help
		err = cli.ShowCommandHelp(ctx, "gadget")
	}
	return err
}
