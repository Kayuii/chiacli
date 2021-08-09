package base

import (
	target "github.com/kayuii/chiacli/mass"
	"github.com/urfave/cli/v2"
)

var massFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "overwrite",
		Aliases: []string{"o"},
		Usage:   "overwrite existed binding list file",
		Value:   true,
	},
	&cli.StringFlag{
		Name:  "keystore",
		Usage: "specify the keystore to eliminate files without private key",
		Value: "",
	},
	&cli.StringSliceFlag{
		Name:    "dirlist",
		Aliases: []string{"d"},
		Usage:   "specify the searching directories",
		Value:   nil,
	},
	&cli.StringFlag{
		Name:    "dirs",
		Aliases: []string{"D"},
		Usage:   "specify the searching directories and subdirectories",
		Value:   "",
	},
}

func NewMassCli() *cli.Command {
	return &cli.Command{
		Name:    "MassBindingTarget",
		Aliases: []string{"mt"},
		Usage:   "Mass Binding Target",
		Action:  target.Target,
		Flags:   massFlags,
	}
}
