package cmds

import (
	"github.com/urfave/cli/v2"

	"github.com/2bit-software/gogo/pkg/funcs"
)

// scriptAction by default just lists the functions
func scriptAction(ctx *cli.Context) error {
	args := ctx.Args().Slice()
	opts, err := BuildOptions(ctx)
	if err != nil {
		return err
	}
	// run the command
	err = scripts.Run(opts, args)
	count, err := scripts.ShowFuncList(opts)
	if err != nil {
		return err
	}
	if count == 0 {
		// show command help
		err = cli.ShowSubcommandHelp(ctx)
		return nil
	}
	return nil
}

func ScriptCommand() *cli.Command {
	return &cli.Command{
		Name:        "script",
		Usage:       "script [subcommand] [arguments...]",
		Description: `Manage scripts. This includes pre-building binaries, performing checks, and listing.`,
		Action:      scriptAction,
		Subcommands: []*cli.Command{
			BuildCommand(),
			// TODO: perform audits/checks?
		},
	}
}
