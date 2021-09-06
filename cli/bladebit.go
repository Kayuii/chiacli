package base

import (
	"github.com/kayuii/chiacli/plot"
	"github.com/urfave/cli/v2"
)

var bladebitFlags = []cli.Flag{
	&cli.IntFlag{
		Name:        NumPlots,
		Aliases:     []string{"n"},
		Value:       1,
		DefaultText: "default = 1, -1 = infinite",
		Usage:       "Number of plots to create. ",
	},
	&cli.IntFlag{
		Name:    KSize,
		Aliases: []string{"k"},
		Value:   32,
		Hidden:  true,
		Usage:   "Plot size. ",
	},
	&cli.IntFlag{
		Name:    Threads,
		Aliases: []string{"r"},
		Value:   4,
		Usage:   "Number of threads. ",
	},
	&cli.BoolFlag{
		Name:    Progress,
		Aliases: []string{"P"},
		Value:   false,
		Usage:   "Display progress percentage during plotting. ",
	},
	&cli.StringFlag{
		Name:    FinalPath,
		Aliases: []string{"d"},
		Value:   ".",
		Usage:   "Final directory. ",
	},
	&cli.StringFlag{
		Name:    LogPath,
		Aliases: []string{"l"},
		Value:   "./logs",
		Hidden:  true,
		Usage:   "Logs directory. ",
	},
	&cli.StringFlag{
		Name:    FarmePublicKey,
		Aliases: []string{"fpk", "f"},
		Value:   "96160804d76ccb56d937536935da2f5ecd32b19d55b56c1ca6c9bc24044ef1d118a8d773ec146130354f19a43483bac0",
		Usage:   "The farmer public key. ",
	},
	&cli.StringFlag{
		Name:    PoolPublicKey,
		Aliases: []string{"ppk", "p"},
		Value:   "b6e26610006b42b33bbc458dc42e8a41bcf25403382dd0074d61679a792f3570e54c22bca6d9863f6c4b22a68355e614",
		Usage:   "The pool public key. ",
	},
	&cli.StringFlag{
		Name:    PoolContractAddress,
		Aliases: []string{"pca", "c"},
		Value:   "",
		Usage:   "The Pool Contract Address",
	},
	&cli.IntFlag{
		Name:    Sleep,
		Aliases: []string{"s"},
		Value:   60,
		Hidden:  true,
		Usage:   "Sleep. ",
	},
}

func bladebitAction(c *cli.Context) error {
	config := &plot.Config{
		NumPlots:            c.Int(NumPlots),
		KSize:               c.Int(KSize),
		Threads:             c.Int(Threads),
		Progress:            c.Bool(Progress),
		FinalPath:           c.String(FinalPath),
		LogPath:             c.String(LogPath),
		FarmePublicKey:      c.String(FarmePublicKey),
		PoolPublicKey:       c.String(PoolPublicKey),
		PoolContractAddress: c.String(PoolContractAddress),
		Sleep:               c.Int(Sleep),
	}
	return plot.New().Bladebit(config)
}

func NewBladebit() *cli.Command {
	return &cli.Command{
		Name:    "FastPos",
		Aliases: []string{"bladebit"},
		Usage:   "from harold-b/bladebit",
		Action:  bladebitAction,
		Flags:   bladebitFlags,
	}
}
