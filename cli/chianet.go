package base

import (
	"github.com/kayuii/chiacli/plot"
	"github.com/urfave/cli/v2"
)

var plotFlags = []cli.Flag{
	&cli.IntFlag{
		Name:    NumPlots,
		Aliases: []string{"n"},
		Value:   1,
		Usage:   "Batch plotting count. ",
	},
	&cli.IntFlag{
		Name:    KSize,
		Aliases: []string{"k"},
		Value:   32,
		Usage:   "Plot size. ",
	},
	&cli.IntFlag{
		Name:    Stripes,
		Aliases: []string{"s"},
		Value:   65536,
		Usage:   "Size of stripes.",
	},
	&cli.IntFlag{
		Name:    Buffer,
		Aliases: []string{"b"},
		Hidden:  true,
		Value:   3389,
		Usage:   "Megabytes to be used as buffer for sorting and plotting. ",
	},
	&cli.IntFlag{
		Name:    Threads,
		Aliases: []string{"r"},
		Value:   2,
		Usage:   "Number of threads. ",
	},
	&cli.IntFlag{
		Name:    Buckets,
		Aliases: []string{"u"},
		Hidden:  true,
		Value:   128,
		Usage:   "Number of buckets. ",
	},
	&cli.BoolFlag{
		Name:    NoBitfield,
		Aliases: []string{"e"},
		Value:   false,
		Usage:   "Disable bitfield. ",
	},
	&cli.BoolFlag{
		Name:    Progress,
		Aliases: []string{"P"},
		Value:   false,
		Usage:   "Display progress percentage during plotting. ",
	},
	&cli.StringFlag{
		Name:    TempPath,
		Aliases: []string{"t"},
		Value:   ".",
		Usage:   "Temporary directory. ",
	},
	&cli.StringFlag{
		Name:    Temp2Path,
		Aliases: []string{"2"},
		Value:   ".",
		Usage:   "Second Temporary directory. ",
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
}

func chiaAction(c *cli.Context) error {
	config := &plot.Config{
		NumPlots:            c.Int(NumPlots),
		KSize:               c.Int(KSize),
		Stripes:             c.Int(Stripes),
		Buffer:              c.Int(Buffer),
		Threads:             c.Int(Threads),
		Buckets:             c.Int(Buckets),
		NoBitfield:          c.Bool(NoBitfield),
		Progress:            c.Bool(Progress),
		TempPath:            c.String(TempPath),
		Temp2Path:           c.String(Temp2Path),
		FinalPath:           c.String(FinalPath),
		LogPath:             c.String(LogPath),
		FarmePublicKey:      c.String(FarmePublicKey),
		PoolPublicKey:       c.String(PoolPublicKey),
		PoolContractAddress: c.String(PoolContractAddress),
	}
	return plot.New().Chia(config)
}

func NewChiaNet() *cli.Command {
	return &cli.Command{
		Name:    "Chia",
		Aliases: []string{"chia"},
		Usage:   "chia-blockchain",
		Action:  chiaAction,
		Flags:   plotFlags,
	}
}
