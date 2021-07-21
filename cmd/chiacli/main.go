package main

import (
	"log"
	"os"
	"time"

	base "github.com/kayuii/chiacli/cli"

	"github.com/kayuii/chiacli"
	"github.com/urfave/cli/v2"
)

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
		base.NewChiaNet(),
		base.NewFastpos(),
		base.NewProofOfSpace(),
		base.NewMassCli(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
