package main

import (
	"bufio"
	"github.com/gavincabbage/stegosaurus"
	"github.com/urfave/cli"
	"io"
	"os"
)

var encodeCmd = cli.Command{
	Name:        "encode",
	ShortName:   "e",
	Usage:       "Encode a signal",
	Description: "Encode a signal",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "payload, p",
			Required: true,
		},
		cli.StringFlag{
			Name:     "carrier, c",
			Required: true,
		},
		cli.StringFlag{
			Name: "outfile, o",
		},
	},
	Action: encode,
}

func encode(ctx *cli.Context) (err error) {
	p, err := reader(ctx.String("payload"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	c, err := reader(ctx.String("carrier"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	var w io.Writer = os.Stdout
	if path := ctx.String("outfile"); path != "" {
		f, err := os.Create(path)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer func() {
			if e := f.Close(); e != nil && err == nil {
				err = cli.NewExitError(e, 1)
			}
		}()

		b := bufio.NewWriter(f)
		defer func() {
			if e := b.Flush(); e != nil && err == nil {
				err = cli.NewExitError(e, 1)
			}
		}()

		w = b
	}

	if err = stegosaurus.Encode(p, c, w); err != nil {
		return cli.NewExitError(err, 1)
	}

	return
}
