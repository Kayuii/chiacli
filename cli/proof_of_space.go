package base

import (
	"github.com/kayuii/chiacli/plot"
	"github.com/urfave/cli/v2"
)

func chiaposAction(c *cli.Context) error {
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
	return plot.New().Pos(config)
}

func NewProofOfSpace() *cli.Command {
	return &cli.Command{
		Name:    "ProofOfSpace",
		Aliases: []string{"pos"},
		Usage:   "Chia Proof of Space",
		Action:  chiaposAction,
		Flags:   plotFlags,
	}
}
