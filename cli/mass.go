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
		Value:   false,
	},
	&cli.BoolFlag{
		Name:    "all",
		Aliases: []string{"a"},
		Usage:   "list all files instead of only plotted files",
		Value:   false,
	},
	&cli.StringFlag{
		Name:  "keystore",
		Usage: "specify the keystore to eliminate files without private key",
		Value: "",
	},
	&cli.StringSliceFlag{
		Name:    "dirs",
		Aliases: []string{"d"},
		Usage:   "specify the searching directories",
		Value:   nil,
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
