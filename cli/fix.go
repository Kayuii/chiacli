package base

import (
	"github.com/kayuii/chiacli/fix"
	"github.com/urfave/cli/v2"
)

var memoFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    FarmePublicKey,
		Aliases: []string{"fpk"},
		Value:   "96160804d76ccb56d937536935da2f5ecd32b19d55b56c1ca6c9bc24044ef1d118a8d773ec146130354f19a43483bac0",
		Usage:   "The farmer public key. ",
	},
	&cli.StringFlag{
		Name:    PoolPublicKey,
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
		FarmePublicKey: c.String(FarmePublicKey),
		PoolPublicKey:  c.String(PoolPublicKey),
		LocalSk:        c.String(LocalSk),
		Memo:           c.String(Memo),
		FilePath:       c.String(FilePath),
		Pattern:        c.String(Pattern),
	}
	return fix.New().Check(config)
}

func NewFix() *cli.Command {
	return &cli.Command{
		Name:    "Check",
		Aliases: []string{"check"},
		Usage:   "check plos",
		Hidden:  true,
		Action:  memoAction,
		Flags:   memoFlags,
	}
}
