package main

import (
	"log"
	"os"
	"time"

	"github.com/kayuii/chiacli"
	"github.com/kayuii/chiacli/fix"
	"github.com/kayuii/chiacli/plot"
	"github.com/urfave/cli/v2"
)

const (
	NumPlots   = "NumPlots"
	KSize      = "KSize"
	Stripes    = "Stripes"
	Buffer     = "Buffer"
	Threads    = "Threads"
	Buckets    = "Buckets"
	NoBitfield = "NoBitfield"
	Progress   = "Progress"
	TempPath   = "TempPath"
	Temp2Path  = "Temp2Path"
	FinalPath  = "FinalPath"
	Total      = "Total"
	Sleep      = "Sleep"
	LogPath    = "LogPath"
	FarmerKey  = "FarmerKey"
	PoolKey    = "PoolKey"
	LocalSk    = "LocalSk"
	Memo       = "Memo"
	FilePath   = "FilePath"
	Pattern    = "Pattern"
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
		Aliases: []string{"p"},
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
		Name:    FarmerKey,
		Aliases: []string{"fpk"},
		Value:   "96160804d76ccb56d937536935da2f5ecd32b19d55b56c1ca6c9bc24044ef1d118a8d773ec146130354f19a43483bac0",
		Usage:   "The farmer public key. ",
	},
	&cli.StringFlag{
		Name:    PoolKey,
		Aliases: []string{"ppk"},
		Value:   "b6e26610006b42b33bbc458dc42e8a41bcf25403382dd0074d61679a792f3570e54c22bca6d9863f6c4b22a68355e614",
		Usage:   "The pool public key. ",
	},
}

func chiaAction(c *cli.Context) error {
	config := &plot.Config{
		NumPlots:   c.Int(NumPlots),
		KSize:      c.Int(KSize),
		Stripes:    c.Int(Stripes),
		Buffer:     c.Int(Buffer),
		Threads:    c.Int(Threads),
		Buckets:    c.Int(Buckets),
		NoBitfield: c.Bool(NoBitfield),
		Progress:   c.Bool(Progress),
		TempPath:   c.String(TempPath),
		Temp2Path:  c.String(Temp2Path),
		FinalPath:  c.String(FinalPath),
		LogPath:    c.String(LogPath),
		FarmerKey:  c.String(FarmerKey),
		PoolKey:    c.String(PoolKey),
	}
	return plot.New().Chia(config)
}

func chiaposAction(c *cli.Context) error {
	config := &plot.Config{
		NumPlots:   c.Int(NumPlots),
		KSize:      c.Int(KSize),
		Stripes:    c.Int(Stripes),
		Buffer:     c.Int(Buffer),
		Threads:    c.Int(Threads),
		Buckets:    c.Int(Buckets),
		NoBitfield: c.Bool(NoBitfield),
		Progress:   c.Bool(Progress),
		TempPath:   c.String(TempPath),
		Temp2Path:  c.String(Temp2Path),
		FinalPath:  c.String(FinalPath),
		LogPath:    c.String(LogPath),
		FarmerKey:  c.String(FarmerKey),
		PoolKey:    c.String(PoolKey),
	}
	return plot.New().Pos(config)
}

var fastplotFlags = []cli.Flag{
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
	&cli.IntFlag{
		Name:    Buckets,
		Aliases: []string{"u"},
		Hidden:  true,
		Value:   256,
		Usage:   "Number of buckets. ",
	},
	&cli.BoolFlag{
		Name:    Progress,
		Aliases: []string{"p"},
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
		Name:    FarmerKey,
		Aliases: []string{"fpk"},
		Value:   "96160804d76ccb56d937536935da2f5ecd32b19d55b56c1ca6c9bc24044ef1d118a8d773ec146130354f19a43483bac0",
		Usage:   "The farmer public key. ",
	},
	&cli.StringFlag{
		Name:    PoolKey,
		Aliases: []string{"ppk"},
		Value:   "b6e26610006b42b33bbc458dc42e8a41bcf25403382dd0074d61679a792f3570e54c22bca6d9863f6c4b22a68355e614",
		Usage:   "The pool public key. ",
	},
}

func fastposAction(c *cli.Context) error {
	config := &plot.Config{
		NumPlots:  c.Int(NumPlots),
		KSize:     c.Int(KSize),
		Threads:   c.Int(Threads),
		Buckets:   c.Int(Buckets),
		Progress:  c.Bool(Progress),
		TempPath:  c.String(TempPath),
		Temp2Path: c.String(Temp2Path),
		FinalPath: c.String(FinalPath),
		LogPath:   c.String(LogPath),
		FarmerKey: c.String(FarmerKey),
		PoolKey:   c.String(PoolKey),
	}
	return plot.New().FastPos(config)
}

var memoFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    FarmerKey,
		Aliases: []string{"fpk"},
		Value:   "96160804d76ccb56d937536935da2f5ecd32b19d55b56c1ca6c9bc24044ef1d118a8d773ec146130354f19a43483bac0",
		Usage:   "The farmer public key. ",
	},
	&cli.StringFlag{
		Name:    PoolKey,
		Aliases: []string{"ppk"},
		Value:   "b6e26610006b42b33bbc458dc42e8a41bcf25403382dd0074d61679a792f3570e54c22bca6d9863f6c4b22a68355e614",
		Usage:   "The pool public key. ",
	},
	&cli.StringFlag{
		Name:    LocalSk,
		Aliases: []string{"sk"},
		Value:   "",
		Usage:   "Local sk. ",
	},
	&cli.StringFlag{
		Name:    Memo,
		Aliases: []string{"m"},
		Value:   "",
		Usage:   "memo. ",
	},
	&cli.StringFlag{
		Name:    FilePath,
		Aliases: []string{"d", "dir"},
		Value:   ".",
		Usage:   "dir. ",
	},
	&cli.StringFlag{
		Name:    Pattern,
		Aliases: []string{"p", "pattern"},
		Value:   `.(plot)$`,
		Usage:   "pattern. ",
	},
}

func memoAction(c *cli.Context) error {
	config := &fix.Config{
		FarmerKey: c.String(FarmerKey),
		PoolKey:   c.String(PoolKey),
		LocalSk:   c.String(LocalSk),
		Memo:      c.String(Memo),
		FilePath:  c.String(FilePath),
		Pattern:   c.String(Pattern),
	}
	return fix.New().Check(config)
}

func main() {

	app := cli.NewApp()
	app.Version = chiacli.Version
	app.Usage = "plotting utility for chia."
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "Kayuii",
			Email: "577738@qq.com",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "Chia",
			Aliases: []string{"chia"},
			Usage:   "chia-blockchain",
			Action:  chiaAction,
			Flags:   plotFlags,
		},
		{
			Name:    "FastPos",
			Aliases: []string{"fastpos"},
			Usage:   "from madMAx43v3r/chia-plotter",
			Action:  fastposAction,
			Flags:   fastplotFlags,
		},
		{
			Name:    "ProofOfSpace",
			Aliases: []string{"pos"},
			Usage:   "Chia Proof of Space",
			Action:  chiaposAction,
			Flags:   plotFlags,
		},
		{
			Name:    "Check",
			Aliases: []string{"check"},
			Usage:   "check plos",
			Action:  memoAction,
			Flags:   memoFlags,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
