package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/gavincabbage/stegosaurus/image"
)

var version = "dev"

func main() {
	app := cli.App{
		Name:        "stegosaurus",
		HelpName:    "stegosaurus",
		Version:     version,
		Description: "Steganography tool",
		Commands: []cli.Command{
			{
				Name:        "encode",
				ShortName:   "e",
				Usage:       "Encode a secret payload in a carrier signal",
				Description: "Encode a secret payload in a carrier signal",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "payload, p",
						Usage:    "Payload; Defaults to stdin; Use @path to read from a file",
						Required: true,
					},
					cli.StringFlag{
						Name:     "carrier, c",
						Usage:    "Carrier file",
						Required: true,
					},
					cli.StringFlag{
						Name:  "outfile, o",
						Usage: "Output file; Defaults to stdout",
					},
					cli.StringFlag{
						Name:  "algorithm, a",
						Usage: "Encoding algorithm; Defaults to LSB",
					},
					cli.StringFlag{
						Name:  "key, k",
						Usage: "Encoding key, if required by algorithm",
					},
				},
				Action: encode,
			},
			{
				Name:        "decode",
				ShortName:   "d",
				Usage:       "Decode a secret payload from a carrier package",
				Description: "Decode a secret payload from a carrier package",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "package, p",
						Usage: "Encoded carrier package file; Defaults to stdin",
					},
					cli.StringFlag{
						Name:  "outfile, o",
						Usage: "Output file; Defaults to stdout",
					},
					cli.StringFlag{
						Name:  "algorithm, a",
						Usage: "Encoding algorithm; Defaults to LSB",
					},
					cli.StringFlag{
						Name:  "key, k",
						Usage: "Encoding key, if required by algorithm",
					},
				},
				Action: decode,
			},
		},
	}

	cli.HandleExitCoder(app.Run(os.Args))
}

func encode(ctx *cli.Context) (err error) {
	defer func() {
		if err != nil {
			err = cli.NewExitError(err, 1)
		}
	}()

	var p io.Reader = os.Stdin
	if payload := ctx.String("payload"); payload != "" {
		if strings.HasPrefix(payload, "@") {
			d, err := ioutil.ReadFile(strings.TrimPrefix(payload, "@"))
			if err != nil {
				return err
			}
			p = bytes.NewReader(d)
		} else {
			p = strings.NewReader(payload)
		}
	}

	b, err := ioutil.ReadFile(ctx.String("carrier"))
	if err != nil {
		return err
	}
	c := bytes.NewReader(b)

	var w io.Writer = os.Stdout
	if path := ctx.String("outfile"); path != "" {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			if e := f.Close(); e != nil && err == nil {
				err = e
			}
		}()

		b := bufio.NewWriter(f)
		defer func() {
			if e := b.Flush(); e != nil && err == nil {
				err = e
			}
		}()

		w = b
	}

	var (
		key     = []byte(ctx.String("key"))
		encoder = image.NewEncoder(key)
	)
	if err = encoder.Encode(p, c, w); err != nil {
		return err
	}

	return nil
}

func decode(ctx *cli.Context) (err error) {
	defer func() {
		if err != nil {
			err = cli.NewExitError(err, 1)
		}
	}()

	var r io.Reader = os.Stdin
	if path := ctx.String("package"); path != "" {
		d, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		r = bytes.NewReader(d)
	}

	var w io.Writer = os.Stdout
	if path := ctx.String("outfile"); path != "" {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			if e := f.Close(); e != nil && err == nil {
				err = e
			}
		}()

		b := bufio.NewWriter(f)
		defer func() {
			if e := b.Flush(); e != nil && err == nil {
				err = e
			}
		}()

		w = b
	}

	var (
		key     = []byte(ctx.String("key"))
		encoder = image.NewEncoder(key)
	)
	if err = encoder.Decode(r, w); err != nil {
		return err
	}

	return nil
}
