package cmds

import (
	"fmt"
	"github.com/2bit-software/gogo/pkg/gadgets"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "init [folder] [--gogomodpath] [--global]",
		Description: `Initialize a new GoGo workspace. Creates a .gogo subfolder either in the current directory or in the specified folder.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "gogomodpath",
				Usage:    "The path to the GoGo module. Defaults to the current directory.",
				Required: false,
				Value:    "",
				Aliases:  []string{"gmp"},
			},
			&cli.BoolFlag{
				Name:     "global",
				Usage:    "Initialize the GoGo workspace globally.",
				Required: false,
				Aliases:  []string{"g"},
			},
		},
		Action: initAction,
	}
}

// initAction sets up the environment to be used by GoGo.
// In most cases, this means making a .gogo folder in the current directory.
// It is recommended that the .GoGo folder have it's own go.mod, so we make one by default.
// If you don't pass a --godmodpath flag, we will use "notgithub.com/gogo/gadgets"
func initAction(context *cli.Context) error {
	// get the folder to initialize
	folder := context.Args().First()

	// get the path to the GoGo module
	gogomodpath := context.String("gogomodpath")
	if gogomodpath == "" {
		gogomodpath = "notgithub.com/gogo/gadgets"
	}

	// if the global flag is set, we need to make the .gogo folder in the home directory
	// unless the folder is specified, then we just use it
	if context.Bool("global") && folder == "" {
		// set folder to the home directory
		hmdir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		folder = path.Join(hmdir, ".gogo")
	}

	if folder == "" {
		folder = path.Join(".", ".gogo")
	}

	// call the gadgets.InitCommand function
	err := gadgets.Init(path.Join(folder), gogomodpath)
	if err != nil {
		return err
	}
	fmt.Printf("Initialized GoGo workspace at %s\n", folder)
	return nil
}
