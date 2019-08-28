package main

import (
	"bufio"
	"github.com/gavincabbage/stegosaurus"
	"github.com/urfave/cli"
	"io"
	"os"
)

var decodeCmd = cli.Command{
	Name:        "decode",
	ShortName:   "d",
	Usage:       "Decode a signal",
	Description: "Decode a signal",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "data, d",
			Required: true,
		},
		cli.StringFlag{
			Name: "outfile, o",
		},
	},
	Action: decode,
}

func decode(ctx *cli.Context) (err error) {
	r, err := reader(ctx.String("data"))
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

	if err = stegosaurus.Decode(r, w); err != nil {
		return cli.NewExitError(err, 1)
	}

	return
}
