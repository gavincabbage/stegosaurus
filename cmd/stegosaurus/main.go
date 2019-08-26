package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var Version string

func main() {
	app := cli.App{
		Name:        "stegosaurus",
		HelpName:    "stegosaurus",
		Version:     Version,
		Description: "stegonography encoding and decoding tool",
		Author:      "Gavin Cabbage",
		Commands: []cli.Command{
			{
				Name:      "encode",
				ShortName: "e",
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
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:      "decode",
				ShortName: "d",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "data, d",
						Required: true,
					},
					cli.StringFlag{
						Name: "outfile, o",
					},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
